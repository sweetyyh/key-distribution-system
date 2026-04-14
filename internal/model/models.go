package model

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleBuyer UserRole = "buyer"
)

type User struct {
	ID           uint64    `gorm:"primaryKey"`
	Username     string    `gorm:"size:64;not null;uniqueIndex"`
	Email        string    `gorm:"size:128;not null;uniqueIndex"`
	PasswordHash string    `gorm:"size:256;not null"`
	Role         UserRole  `gorm:"type:enum('admin','buyer');not null;default:'buyer'"`
	Status       int8      `gorm:"not null;default:1"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type Category struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"size:64;not null"`
	ParentID  uint64 `gorm:"not null;default:0"`
	SortOrder int    `gorm:"not null;default:0"`
}

type Product struct {
	ID             uint64    `gorm:"primaryKey"`
	CategoryID     uint64    `gorm:"not null;index"`
	Name           string    `gorm:"size:128;not null"`
	Description    string    `gorm:"type:text"`
	Price          string    `gorm:"type:decimal(18,2);not null"`
	WholesaleRules string    `gorm:"type:json"`
	Stock          int       `gorm:"not null;default:0"`
	Status         int8      `gorm:"not null;default:1"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

type CardKey struct {
	ID        uint64  `gorm:"primaryKey"`
	ProductID uint64  `gorm:"not null;index:idx_product_status_id,priority:1"`
	Content   string  `gorm:"type:text;not null"`
	Status    int8    `gorm:"not null;default:0;index:idx_product_status_id,priority:2"`
	OrderID   *uint64 `gorm:"index"`
	LockedAt  *time.Time
	SoldAt    *time.Time
}

type Order struct {
	ID          uint64    `gorm:"primaryKey"`
	OrderNo     string    `gorm:"size:32;not null;uniqueIndex"`
	UserID      uint64    `gorm:"not null;index"`
	ProductID   uint64    `gorm:"not null;index"`
	Quantity    int       `gorm:"not null"`
	UnitPrice   string    `gorm:"type:decimal(18,2);not null"`
	TotalAmount string    `gorm:"type:decimal(18,2);not null"`
	Status      int8      `gorm:"not null;default:0;index"`
	PayChannel  string    `gorm:"size:32;not null"`
	ExpiresAt   time.Time `gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	PaidAt      *time.Time
}

type Payment struct {
	ID          uint64 `gorm:"primaryKey"`
	OrderID     uint64 `gorm:"not null;index"`
	OrderNo     string `gorm:"size:32;not null;index"`
	OutTradeNo  string `gorm:"size:64;not null;uniqueIndex"`
	TradeNo     string `gorm:"size:64;index:idx_channel_trade,priority:2"`
	Channel     string `gorm:"size:32;not null;index:idx_channel_trade,priority:1"`
	Amount      string `gorm:"type:decimal(18,2);not null"`
	Status      int8   `gorm:"not null;default:0;index"`
	RawCallback string `gorm:"type:json"`
	PaidAt      *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type OrderItem struct {
	ID              uint64    `gorm:"primaryKey"`
	OrderID         uint64    `gorm:"not null;index"`
	CardKeyID       uint64    `gorm:"not null;index"`
	ContentSnapshot string    `gorm:"type:text;not null"`
	DeliveredAt     time.Time `gorm:"not null"`
}

type ProductSnapshot struct {
	ID             uint64    `gorm:"primaryKey"`
	OrderID        uint64    `gorm:"not null;index"`
	ProductID      uint64    `gorm:"not null;index"`
	Name           string    `gorm:"size:128;not null"`
	Price          string    `gorm:"type:decimal(18,2);not null"`
	WholesaleRules string    `gorm:"type:json"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
