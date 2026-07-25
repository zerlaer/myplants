package config

import (
	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Reminder ReminderConfig `mapstructure:"reminder"`
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

var cfg *Config

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	cfg = c
	return c, nil
}

// Get 获取配置
func Get() *Config {
	return cfg
}
