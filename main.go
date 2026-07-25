package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"myplants/internal/config"
	"myplants/internal/database"
	"myplants/internal/logger"
	"myplants/internal/router"
)

func main() {
	configPath := flag.String("c", "", "配置文件路径")
	flag.Parse()

	if *configPath != "" {
		os.Setenv("VIPER_CONFIG", *configPath)
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Logger); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.L().Sync()

	logger.S().Infof("我的花园 启动中... port=%d mode=%s", cfg.Server.Port, cfg.Server.Mode)

	// 初始化上传目录
	if err := os.MkdirAll(cfg.Upload.Path, 0755); err != nil {
		logger.S().Errorf("创建上传目录失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.S().Fatalf("数据库初始化失败: %v", err)
	}

	// 启动 HTTP 服务
	r := router.Setup(cfg)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		logger.S().Infof("HTTP 服务已启动 addr=%s", addr)
		if err := r.Run(addr); err != nil {
			logger.S().Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.S().Info("正在关闭服务...")
	logger.S().Info("服务已关闭")
}
