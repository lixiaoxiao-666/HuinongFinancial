package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/api"
	"backend/internal/conf"
	"backend/internal/data"
	"backend/pkg"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	config, err := conf.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger, err := initLogger(config)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("应用启动中...",
		zap.String("name", config.App.Name),
		zap.String("version", config.App.Version),
		zap.String("env", config.App.Env),
		zap.Int("port", config.Server.Port),
	)

	// 初始化数据层
	dataLayer, err := data.NewData(config, logger)
	if err != nil {
		logger.Fatal("初始化数据层失败", zap.Error(err))
	}
	defer dataLayer.Close()

	// 自动迁移数据库表结构
	if err := dataLayer.AutoMigrate(); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("数据库迁移完成")

	// 初始化JWT管理器
	jwtManager := pkg.NewJWTManager(config.JWT.Secret, config.JWT.Expire)

	// 初始化路由器
	router := api.NewRouter(config, dataLayer, jwtManager, logger)
	engine := router.SetupRoutes()

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", config.Server.Port)
	logger.Info("服务器启动成功", zap.String("addr", serverAddr))

	// 优雅关闭
	go func() {
		if err := engine.Run(serverAddr); err != nil {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")
	logger.Info("服务器已关闭")
}

// initLogger 初始化日志
func initLogger(config *conf.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if config.App.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
