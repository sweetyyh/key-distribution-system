package store

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"key-distribution-system/internal/model"
)

// RegisterInput 注册输入
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// Register 注册新用户，返回创建的用户
func (s *DBStore) Register(input RegisterInput) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         model.RoleBuyer,
		Status:       1,
	}

	if err = s.db.Create(user).Error; err != nil {
		// 唯一约束冲突
		return nil, fmt.Errorf("username or email already exists")
	}
	return user, nil
}

// LoginByUsername 通过用户名和密码登录，返回用户
func (s *DBStore) LoginByUsername(username, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ? AND status = 1", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &user, nil
}

// GetUserByID 通过 ID 查询用户
func (s *DBStore) GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
