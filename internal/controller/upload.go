package controller

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // 注册 webp 解码器

	"myplants/internal/config"
)

// ServeUpload 提供 /uploads 静态图片访问:
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
	thumb, err := resizeToJPEG(src, w, h)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(thumbPath), 0755); err != nil {
		return "", err
	}
	tmp := thumbPath + ".tmp"
	if err := os.WriteFile(tmp, thumb, 0644); err != nil {
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

func resizeToJPEG(src image.Image, w, h int) ([]byte, error) {
	srcBounds := src.Bounds()

	// 透明像素铺白底,避免 PNG 转 JPEG 后变黑
	flat := image.NewNRGBA(image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy()))
	draw.Draw(flat, flat.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(flat, flat.Bounds(), src, srcBounds.Min, draw.Over)

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), flat, flat.Bounds(), draw.Src, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 82}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
