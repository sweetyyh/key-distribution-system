package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"key-distribution-system/internal/db"
	"key-distribution-system/internal/model"
)

func main() {
	_ = godotenv.Load(".env")

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("请在 .env 中设置 ADMIN_USERNAME 和 ADMIN_PASSWORD")
	}

	if err := db.Init("data/kds.db"); err != nil {
		log.Fatalf("init db: %v", err)
	}

	// 检查是否已存在
	var existing model.User
	if err := db.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		log.Fatalf("用户 %s 已存在", username)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	user := model.User{
		Username:     username,
		Email:        username + "@admin.local",
		PasswordHash: string(hash),
		Role:         model.RoleAdmin,
		Status:       1,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Fatalf("create admin: %v", err)
	}

	fmt.Printf("管理员创建成功: %s\n", username)
}
