package controller

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"myplants/internal/config"
	"myplants/internal/imgutil"
	"myplants/internal/storage"
)

// ServeUpload 提供 /uploads 静态图片访问(本地驱动):
// 支持 ?w=<宽度> 按需生成缩略图(磁盘缓存在 uploads/.thumbs),并为图片添加长缓存头
func ServeUpload(c *gin.Context) {
	cfg := config.Get()
	uploadRoot := filepath.Clean(cfg.Upload.Path)

	rel := filepath.ToSlash(filepath.Clean("/" + strings.TrimPrefix(c.Param("filepath"), "/")))
	full := filepath.Join(uploadRoot, filepath.FromSlash(rel))
	if !strings.HasPrefix(full, uploadRoot+string(os.PathSeparator)) {
		c.Status(http.StatusForbidden)
		return
	}

	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		c.Status(http.StatusNotFound)
		return
	}

	if wStr := c.Query("w"); wStr != "" && !strings.HasSuffix(rel, ".gif") {
		if w, cerr := strconv.Atoi(wStr); cerr == nil && w > 0 && w <= 4096 {
			if thumbPath, terr := ensureThumb(full, rel, w, info.Size()); terr == nil && thumbPath != "" {
				full = thumbPath
			}
		}
	}

	if full == filepath.Join(uploadRoot, filepath.FromSlash(rel)) {
		// 原图文件名含纳秒时间戳,内容不可变
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		c.Header("Cache-Control", "public, max-age=2592000")
	}
	http.ServeFile(c.Writer, c.Request, full)
}

// ensureThumb 返回缩略图路径;原图窄于请求宽度或处理失败时返回空
func ensureThumb(full, rel string, w int, origSize int64) (string, error) {
	f, err := os.Open(full)
	if err != nil {
		return "", err
	}
	defer f.Close()

	srcCfg, _, err := image.DecodeConfig(f)
	if err != nil || srcCfg.Width <= w {
		return "", err
	}

	thumbName := fmt.Sprintf("%d_%d_%s.jpg", w, origSize, filepath.Base(rel))
	thumbPath := filepath.Join(config.Get().Upload.Path, ".thumbs", filepath.Dir(strings.TrimPrefix(rel, "/")), thumbName)
	if _, err := os.Stat(thumbPath); err == nil {
		return thumbPath, nil
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	src, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	h := src.Bounds().Dy() * w / src.Bounds().Dx()
	data, err := imgutil.EncodeJPEG(imgutil.Resample(src, w, h), 82)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(thumbPath), 0755); err != nil {
		return "", err
	}
	tmp := thumbPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, thumbPath); err != nil {
		// 并发场景下另一个请求已生成同名文件
		os.Remove(tmp)
		if _, serr := os.Stat(thumbPath); serr == nil {
			return thumbPath, nil
		}
		return "", err
	}
	return thumbPath, nil
}

// saveImageFile 保存上传文件:本地驱动写磁盘;R2 驱动上传原图并预生成各档缩略图
func saveImageFile(file *multipart.FileHeader, relPath string) error {
	st := storage.Get()
	key := storage.KeyFromPath(relPath)

	if st.Driver() == "local" {
		full := filepath.Join(config.Get().Upload.Path, filepath.FromSlash(key))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.Create(full)
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	if err := st.Put(key, data, storage.ContentType(key)); err != nil {
		return err
	}
	for w, thumb := range imgutil.Thumbnails(data) {
		if err := st.Put(fmt.Sprintf("w%d/%s", w, key), thumb, "image/jpeg"); err != nil {
			return err
		}
	}
	return nil
}

// deleteImageFile 删除对象;R2 模式同时清掉各档缩略图(尽力而为)
func deleteImageFile(relPath string) {
	st := storage.Get()
	key := storage.KeyFromPath(relPath)
	st.Delete(key)
	if st.Driver() != "local" {
		for _, tk := range storage.TierKeys(key, imgutil.TierWidths) {
			st.Delete(tk)
		}
	}
}
