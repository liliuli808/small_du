package database

import (
	"fmt"
	"time"

	"recipe-ai-backend/internal/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.DBName == "recipe_ai" { // 开发环境
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	sqlDB.SetMaxIdleConns(cfg.MaxConns / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return db, nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
