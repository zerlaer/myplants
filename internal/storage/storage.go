package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"myplants/internal/config"
)

// ErrNotConfigured 驱动配置不完整
var ErrNotConfigured = errors.New("存储驱动未配置完整")

// Storage 图片对象存储接口,key 为 uploads 下的相对路径,如 plant_3/xxx.jpg
type Storage interface {
	Put(key string, data []byte, contentType string) error
	Delete(key string) error
	Driver() string
}

// New 按配置返回存储驱动
func New(cfg *config.Config) (Storage, error) {
	switch cfg.Storage.Driver {
	case "r2":
		if !cfg.Storage.R2Ready() {
			return nil, ErrNotConfigured
		}
		return NewR2(cfg.Storage.R2), nil
	default:
		return &local{root: cfg.Upload.Path}, nil
	}
}

var inst Storage

// Init 在启动时按配置初始化全局存储驱动
func Init(cfg *config.Config) error {
	s, err := New(cfg)
	if err != nil {
		return err
	}
	inst = s
	return nil
}

// Get 返回当前存储驱动,未初始化时回退本地磁盘
func Get() Storage {
	if inst == nil {
		return &local{root: config.Get().Upload.Path}
	}
	return inst
}

// TierKeys 返回某对象在 R2 模式下的全部缩略图 key
func TierKeys(key string, widths []int) []string {
	out := make([]string, 0, len(widths))
	for _, w := range widths {
		out = append(out, "w"+strconv.Itoa(w)+"/"+key)
	}
	return out
}

type local struct {
	root string
}

func (l *local) Driver() string { return "local" }

func (l *local) Put(key string, data []byte, _ string) error {
	full, err := l.resolve(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, data, 0644)
}

func (l *local) Delete(key string) error {
	full, err := l.resolve(key)
	if err != nil {
		return err
	}
	err = os.Remove(full)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (l *local) resolve(key string) (string, error) {
	full := filepath.Join(l.root, filepath.FromSlash(filepath.Clean("/"+strings.TrimPrefix(key, "/"))))
	if !strings.HasPrefix(full, filepath.Clean(l.root)+string(os.PathSeparator)) {
		return "", errors.New("非法路径")
	}
	return full, nil
}

// KeyFromPath 把数据库里的 /uploads/xxx 相对路径转为存储 key
func KeyFromPath(path string) string {
	return strings.TrimPrefix(path, "/uploads/")
}

// ContentType 按扩展名推断
func ContentType(key string) string {
	switch strings.ToLower(filepath.Ext(key)) {
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
