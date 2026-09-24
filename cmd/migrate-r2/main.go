package main

// migrate-r2: 把本地 uploads/ 里的历史图片迁移到 Cloudflare R2
// 用法(在项目根目录):
//   go run ./cmd/migrate-r2 -dry-run          # 只检查配置并列出将上传的对象
//   go run ./cmd/migrate-r2                   # 执行迁移(可重复跑,已存在的会覆盖)
//   go run ./cmd/migrate-r2 -clean            # 删除 R2 bucket 内全部对象(含缩略图)
import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"myplants/internal/config"
	"myplants/internal/imgutil"
	"myplants/internal/storage"
)

// listDeleter 由 r2 客户端实现
type listDeleter interface {
	List(prefix string) ([]string, error)
	Delete(key string) error
}

func main() {
	dryRun := flag.Bool("dry-run", false, "只检查不执行")
	clean := flag.Bool("clean", false, "删除 R2 内全部对象")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v\n", err)
	}
	if !cfg.Storage.R2Ready() {
		fmt.Println("R2 配置不完整。请设置环境变量 STORAGE_DRIVER=r2 以及 R2_ACCOUNT_ID / R2_ACCESS_KEY_ID / R2_SECRET_ACCESS_KEY / R2_BUCKET")
		os.Exit(1)
	}
	st, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("初始化 R2 客户端失败: %v\n", err)
	}

	if *clean {
		ld, ok := st.(listDeleter)
		if !ok {
			log.Fatalf("当前 driver=%s 不支持清理,请确认 STORAGE_DRIVER=r2", st.Driver())
		}
		keys, err := ld.List("")
		if err != nil {
			log.Fatalf("列举 R2 对象失败: %v\n", err)
		}
		fmt.Printf("R2 bucket=%s 共 %d 个对象\n", cfg.Storage.R2.Bucket, len(keys))
		if *dryRun {
			for _, k := range keys {
				fmt.Printf("  将删除: %s\n", k)
			}
			return
		}
		ok2, bad := 0, 0
		for i, k := range keys {
			if err := ld.Delete(k); err != nil {
				fmt.Printf("[%d/%d] 删除失败 %s: %v\n", i+1, len(keys), k, err)
				bad++
				continue
			}
			ok2++
			fmt.Printf("[%d/%d] 已删除 %s\n", i+1, len(keys), k)
		}
		fmt.Printf("清理完成: 成功 %d, 失败 %d\n", ok2, bad)
		if bad > 0 {
			os.Exit(1)
		}
		return
	}

	var files []string
	err = filepath.WalkDir(cfg.Upload.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// .thumbs 是本地按需缓存,R2 用预生成缩略图,无需迁移
			if strings.HasPrefix(d.Name(), ".") && path != cfg.Upload.Path {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatalf("扫描目录失败: %v\n", err)
	}

	fmt.Printf("共发现 %d 个待迁移文件 (driver=r2, bucket=%s)\n", len(files), cfg.Storage.R2.Bucket)
	if *dryRun {
		for _, f := range files {
			key := relKey(cfg, f)
			fmt.Printf("  %s -> %s (+ %dx 缩略图)\n", f, key, len(imgutil.TierWidths))
		}
		return
	}

	ok, failed := 0, 0
	for i, f := range files {
		key := relKey(cfg, f)
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("[%d/%d] 读取失败 %s: %v\n", i+1, len(files), f, err)
			failed++
			continue
		}
		if err := st.Put(key, data, storage.ContentType(key)); err != nil {
			fmt.Printf("[%d/%d] 上传失败 %s: %v\n", i+1, len(files), key, err)
			failed++
			continue
		}
		for w, thumb := range imgutil.Thumbnails(data) {
			if err := st.Put(fmt.Sprintf("w%d/%s", w, key), thumb, "image/jpeg"); err != nil {
				fmt.Printf("[%d/%d] 缩略图上传失败 w%d/%s: %v\n", i+1, len(files), w, key, err)
				failed++
			}
		}
		ok++
		fmt.Printf("[%d/%d] %s (%.1f KB)\n", i+1, len(files), key, float64(len(data))/1024)
	}
	fmt.Printf("完成: 成功 %d, 失败 %d\n", ok, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func relKey(cfg *config.Config, path string) string {
	rel, _ := filepath.Rel(cfg.Upload.Path, path)
	return filepath.ToSlash(rel)
}
