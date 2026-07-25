package database

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"myplants/internal/config"
	"myplants/internal/logger"
	"myplants/internal/model"
)

var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *config.DatabaseConfig) error {
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return err
	}
	DB = db

	if err := autoMigrate(db); err != nil {
		return err
	}

	logger.L().Info("数据库初始化成功", zap.String("path", cfg.Path))
	return nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Plant{},
		&model.Photo{},
		&model.CareRecord{},
		&model.Pot{},
		&model.Repotting{},
		&model.Note{},
	)
}
