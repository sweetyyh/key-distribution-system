package store

import (
	"fmt"
	"key-distribution-system/internal/model"
	"time"
)

// AdminStatsDTO 后台统计数据
type AdminStatsDTO struct {
	TotalRevenue   string  `json:"total_revenue"`
	MonthRevenue   string  `json:"month_revenue"`
	TotalOrders    int64   `json:"total_orders"`
	MonthOrders    int64   `json:"month_orders"`
	TotalUsers     int64   `json:"total_users"`
	TotalStock     int64   `json:"total_stock"`
	TodayOrders    int64   `json:"today_orders"`
	TodayRevenue   string  `json:"today_revenue"`
	PendingOrders  int64   `json:"pending_orders"`
	LowStockCount  int64   `json:"low_stock_count"`
}

// AdminOrderItem 后台订单列表条目
type AdminOrderItem struct {
	OrderNo     string     `json:"order_no"`
	ProductName string     `json:"product_name"`
	Email       string     `json:"email"`
	Quantity    int        `json:"quantity"`
	TotalAmount string     `json:"total_amount"`
	Status      int8       `json:"status"`
	PayChannel  string     `json:"pay_channel"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

// AdminUserItem 后台用户列表条目
type AdminUserItem struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAdminStats 获取后台统计数据
func (s *DBStore) GetAdminStats() *AdminStatsDTO {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	dto := &AdminStatsDTO{}

	// 总订单数 (已支付)
	s.db.Model(&model.Order{}).Where("status = 2").Count(&dto.TotalOrders)

	// 本月订单数
	s.db.Model(&model.Order{}).Where("status = 2 AND created_at >= ?", monthStart).Count(&dto.MonthOrders)

	// 今日订单数
	s.db.Model(&model.Order{}).Where("created_at >= ?", todayStart).Count(&dto.TodayOrders)

	// 待处理订单数 (status=0, not expired)
	s.db.Model(&model.Order{}).Where("status = 0 AND expires_at > ?", now).Count(&dto.PendingOrders)

	// 总用户数
	s.db.Model(&model.User{}).Where("role = 'buyer'").Count(&dto.TotalUsers)

	// 总库存 (available cards)
	s.db.Model(&model.CardKey{}).Where("status = 0").Count(&dto.TotalStock)

	// 低库存商品数 (stock < 50)
	var products []model.Product
	s.db.Where("status = 1").Find(&products)
	for _, p := range products {
		var stock int64
		s.db.Model(&model.CardKey{}).Where("product_id = ? AND status = 0", p.ID).Count(&stock)
		if stock < 50 {
			dto.LowStockCount++
		}
	}

	// 计算营收（只计算 paid 订单的 total_amount sum via raw query）
	type sumResult struct {
		Total string
	}
	var totalRes struct{ Total *string }
	s.db.Model(&model.Order{}).Select("CAST(SUM(CAST(total_amount AS REAL)) AS TEXT) as total").
		Where("status = 2").Scan(&totalRes)
	if totalRes.Total != nil {
		dto.TotalRevenue = formatMoney(*totalRes.Total)
	} else {
		dto.TotalRevenue = "0.00"
	}

	var monthRes struct{ Total *string }
	s.db.Model(&model.Order{}).Select("CAST(SUM(CAST(total_amount AS REAL)) AS TEXT) as total").
		Where("status = 2 AND created_at >= ?", monthStart).Scan(&monthRes)
	if monthRes.Total != nil {
		dto.MonthRevenue = formatMoney(*monthRes.Total)
	} else {
		dto.MonthRevenue = "0.00"
	}

	var todayRes struct{ Total *string }
	s.db.Model(&model.Order{}).Select("CAST(SUM(CAST(total_amount AS REAL)) AS TEXT) as total").
		Where("status = 2 AND created_at >= ?", todayStart).Scan(&todayRes)
	if todayRes.Total != nil {
		dto.TodayRevenue = formatMoney(*todayRes.Total)
	} else {
		dto.TodayRevenue = "0.00"
	}

	return dto
}

func formatMoney(s string) string {
	if s == "" || s == "<nil>" {
		return "0.00"
	}
	return s
}

// AdminListOrders 后台订单列表（最新排序）
func (s *DBStore) AdminListOrders(page, pageSize int, status int) ([]AdminOrderItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	db := s.db.Model(&model.Order{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var orders []model.Order
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminOrderItem, 0, len(orders))
	for _, o := range orders {
		productName := productNameFromSnapshot(s, o.ID, o.ProductID)
		items = append(items, AdminOrderItem{
			OrderNo:     o.OrderNo,
			ProductName: productName,
			Email:       o.Email,
			Quantity:    o.Quantity,
			TotalAmount: o.TotalAmount,
			Status:      o.Status,
			PayChannel:  o.PayChannel,
			CreatedAt:   o.CreatedAt,
			PaidAt:      o.PaidAt,
		})
	}
	return items, total, nil
}

// AdminListUsers 后台用户列表
func (s *DBStore) AdminListUsers(page, pageSize int) ([]AdminUserItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	s.db.Model(&model.User{}).Count(&total)

	var users []model.User
	if err := s.db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	items := make([]AdminUserItem, 0, len(users))
	for _, u := range users {
		items = append(items, AdminUserItem{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      string(u.Role),
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		})
	}
	return items, total, nil
}

// GuestGetOrderCards 无鉴权获取订单卡密（通过邮箱验证）
func (s *DBStore) GuestGetOrderCards(orderNo, email string) (*OrderCardsDTO, error) {
	var order model.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}

	// 如果订单有绑定邮箱，则验证
	if order.Email != "" && order.Email != email {
		return nil, fmt.Errorf("email mismatch")
	}

	if order.Status != 2 {
		return nil, fmt.Errorf("order not paid yet")
	}

	var p model.Product
	productName := ""
	if err := s.db.First(&p, order.ProductID).Error; err == nil {
		productName = p.Name
	}

	var items []model.OrderItem
	s.db.Where("order_id = ?", order.ID).Order("id asc").Find(&items)

	cards := make([]OrderCardItem, 0, len(items))
	for i, item := range items {
		code, pin := splitContent(item.ContentSnapshot)
		cards = append(cards, OrderCardItem{
			Index: i + 1,
			Code:  code,
			Pin:   pin,
		})
	}

	return &OrderCardsDTO{
		OrderNo:     order.OrderNo,
		ProductName: productName,
		Quantity:    order.Quantity,
		TotalAmount: order.TotalAmount,
		PaidAt:      order.PaidAt,
		Cards:       cards,
	}, nil
}
