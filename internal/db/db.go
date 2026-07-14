package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"key-distribution-system/internal/model"
)

var DB *gorm.DB

func Init(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create db dir: %w", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	// 开启 WAL 模式，提升并发读性能
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA foreign_keys=ON")

	// 自动建表
	if err = DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.CardKey{},
		&model.Order{},
		&model.Payment{},
		&model.OrderItem{},
		&model.ProductSnapshot{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}
