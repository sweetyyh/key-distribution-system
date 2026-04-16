package store

import (
	"fmt"
	"time"

	"key-distribution-system/internal/model"
)

// OrderDetailDTO 订单详情
type OrderDetailDTO struct {
	OrderNo     string     `json:"order_no"`
	ProductName string     `json:"product_name"`
	Quantity    int        `json:"quantity"`
	UnitPrice   string     `json:"unit_price"`
	TotalAmount string     `json:"total_amount"`
	Status      int8       `json:"status"`
	PayChannel  string     `json:"pay_channel"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

// OrderListItem 订单列表条目
type OrderListItem struct {
	OrderNo     string     `json:"order_no"`
	ProductName string     `json:"product_name"`
	Quantity    int        `json:"quantity"`
	TotalAmount string     `json:"total_amount"`
	Status      int8       `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

// GetOrderDetail 获取订单详情（校验归属）
func (s *DBStore) GetOrderDetail(orderNo string, userID uint64) (*OrderDetailDTO, error) {
	var order model.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}

	// 从快照取商品名，快照不存在则回退到商品表
	productName := productNameFromSnapshot(s, order.ID, order.ProductID)

	return &OrderDetailDTO{
		OrderNo:     order.OrderNo,
		ProductName: productName,
		Quantity:    order.Quantity,
		UnitPrice:   order.UnitPrice,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		PayChannel:  order.PayChannel,
		ExpiresAt:   order.ExpiresAt,
		CreatedAt:   order.CreatedAt,
		PaidAt:      order.PaidAt,
	}, nil
}

// ListOrders 查询用户订单列表（按创建时间倒序）
func (s *DBStore) ListOrders(userID uint64, page, pageSize int) ([]OrderListItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}

	var total int64
	s.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total)

	var orders []model.Order
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	items := make([]OrderListItem, 0, len(orders))
	for _, o := range orders {
		productName := productNameFromSnapshot(s, o.ID, o.ProductID)
		items = append(items, OrderListItem{
			OrderNo:     o.OrderNo,
			ProductName: productName,
			Quantity:    o.Quantity,
			TotalAmount: o.TotalAmount,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
			PaidAt:      o.PaidAt,
		})
	}
	return items, total, nil
}

func productNameFromSnapshot(s *DBStore, orderID, productID uint64) string {
	var snap model.ProductSnapshot
	if err := s.db.Where("order_id = ?", orderID).First(&snap).Error; err == nil {
		return snap.Name
	}
	var p model.Product
	if err := s.db.Select("name").First(&p, productID).Error; err == nil {
		return p.Name
	}
	return ""
}
