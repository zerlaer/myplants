package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Reminder ReminderConfig `mapstructure:"reminder"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Path   string `mapstructure:"path"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
	Path     string `mapstructure:"path"`
	Filename string `mapstructure:"filename"`
}

type UploadConfig struct {
	Path    string `mapstructure:"path"`
	MaxSize int64  `mapstructure:"max_size"`
}

type ReminderConfig struct {
	DefaultWaterDays     int `mapstructure:"default_water_days"`
	DefaultFertilizeDays int `mapstructure:"default_fertilize_days"`
	DefaultSprayDays     int `mapstructure:"default_spray_days"`
}

// StorageConfig 图片存储驱动: local 本地磁盘 / r2 Cloudflare R2
type StorageConfig struct {
	Driver string   `mapstructure:"driver"`
	R2     R2Config `mapstructure:"r2"`
}

type R2Config struct {
	AccountID       string `mapstructure:"account_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	PublicBase      string `mapstructure:"public_base"` // 公开桶域名,如 https://pub-xxx.r2.dev
}

// R2Ready 判断 R2 配置是否完整可用
func (s StorageConfig) R2Ready() bool {
	return s.Driver == "r2" && s.R2.AccountID != "" && s.R2.AccessKeyID != "" &&
		s.R2.SecretAccessKey != "" && s.R2.Bucket != ""
}

var cfg *Config

// Load 加载配置
func Load() (*Config, error) {
	loadDotEnv(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")
	bindEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	if c.Storage.Driver == "" {
		c.Storage.Driver = "local"
	}
	cfg = c
	return c, nil
}

// loadDotEnv 读取项目根目录 .env 作为环境变量兜底(已存在的环境变量优先,便于容器注入)
func loadDotEnv(dir string) {
	f, err := os.Open(filepath.Join(dir, ".env"))
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}

// 支持用环境变量注入 R2 凭证,避免把密钥写进配置文件
func bindEnv() {
	envs := map[string]string{
		"storage.driver":               "STORAGE_DRIVER",
		"storage.r2.account_id":        "R2_ACCOUNT_ID",
		"storage.r2.access_key_id":     "R2_ACCESS_KEY_ID",
		"storage.r2.secret_access_key": "R2_SECRET_ACCESS_KEY",
		"storage.r2.bucket":            "R2_BUCKET",
		"storage.r2.public_base":       "R2_PUBLIC_BASE",
	}
	for key, env := range envs {
		viper.BindEnv(key, env)
	}
}

// Get 获取配置
func Get() *Config {
	return cfg
}
