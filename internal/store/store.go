package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"key-distribution-system/internal/model"
)

// WholesaleRule 阶梯定价规则
type WholesaleRule struct {
	Min   int     `json:"min"`
	Max   int     `json:"max,omitempty"` // 0 表示无上限
	Price float64 `json:"price"`
}

// ProductDTO 产品数据传输对象
type ProductDTO struct {
	ID             uint64          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Price          string          `json:"price"`
	WholesaleRules []WholesaleRule `json:"wholesale_rules"`
	Stock          int             `json:"stock"`
	Status         int8            `json:"status"`
}

// OrderCardItem 卡密条目
type OrderCardItem struct {
	Index int    `json:"index"`
	Code  string `json:"code"`
	Pin   string `json:"pin"`
}

// OrderCardsDTO 订单卡密响应
type OrderCardsDTO struct {
	OrderNo     string          `json:"order_no"`
	ProductName string          `json:"product_name"`
	Quantity    int             `json:"quantity"`
	TotalAmount string          `json:"total_amount"`
	PaidAt      *time.Time      `json:"paid_at,omitempty"`
	Cards       []OrderCardItem `json:"cards"`
}

// CreateOrderInput 创建订单输入
type CreateOrderInput struct {
	UserID     uint64
	ProductID  uint64
	Quantity   int
	PayChannel string
	Email      string
	Name       string
}

// CreateOrderOutput 创建订单输出
type CreateOrderOutput struct {
	OrderNo     string    `json:"order_no"`
	TotalAmount string    `json:"total_amount"`
	ExpiresAt   time.Time `json:"expires_at"`
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// DBStore 基于 GORM/SQLite 的存储实现
type DBStore struct {
	db *gorm.DB
}

// Global 全局存储实例
var Global *DBStore

// Init 初始化 store 并写入 mock 数据（幂等）
func Init(db *gorm.DB) {
	Global = &DBStore{db: db}
	seedData(db)
}

// CalcPrice 根据阶梯规则和数量计算单价
func CalcPrice(rules []WholesaleRule, qty int) float64 {
	for _, r := range rules {
		if qty >= r.Min && (r.Max == 0 || qty <= r.Max) {
			return r.Price
		}
	}
	if len(rules) > 0 {
		return rules[len(rules)-1].Price
	}
	return 0
}

// generateCardContent 生成卡密内容，格式: SL-XXXX-XXXX-XXXX|DDDD
func generateCardContent() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	segment := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = chars[rng.Intn(len(chars))]
		}
		return string(b)
	}
	code := fmt.Sprintf("SL-%s-%s-%s", segment(4), segment(4), segment(4))
	pin := fmt.Sprintf("%04d", rng.Intn(10000))
	return code + "|" + pin
}

// splitContent 从 "code|pin" 格式中分离 code 和 pin
func splitContent(content string) (string, string) {
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '|' {
			return content[:i], content[i+1:]
		}
	}
	return content, ""
}

// seedData 写入初始商品和卡密（幂等：已有数据则跳过）
func seedData(db *gorm.DB) {
	type productSeed struct {
		id          uint64
		name        string
		description string
		rules       []WholesaleRule
		cardCount   int
	}

	seeds := []productSeed{
		{
			id:          1,
			name:        "Apple App Store",
			description: "Apple App Store gift cards can be used to purchase apps, games, subscriptions, and other content. Valid for US and Global regions.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 5.00}, {Min: 11, Max: 50, Price: 4.90}, {Min: 51, Price: 4.75}},
			cardCount:   200,
		},
		{
			id:          2,
			name:        "Google Play",
			description: "Google Play gift cards let you add funds to your Google Play balance for apps, games, movies, books, and more. Multi-region support.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 10.00}, {Min: 11, Max: 50, Price: 9.60}, {Min: 51, Price: 9.20}},
			cardCount:   200,
		},
		{
			id:          3,
			name:        "Steam Gift Card",
			description: "Steam Gift Cards work just like a gift certificate, which can be redeemed on Steam for the purchase of games, software, wallet credit, and any other item you can purchase on Steam.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 50.00}, {Min: 11, Max: 50, Price: 48.00}, {Min: 51, Price: 45.00}},
			cardCount:   500,
		},
		{
			id:          4,
			name:        "Netflix Prepaid",
			description: "Netflix prepaid gift cards allow users to subscribe or extend their Netflix membership without a credit card. Supports EMEA and US regions.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 16.00}, {Min: 11, Max: 50, Price: 15.50}, {Min: 51, Price: 14.99}},
			cardCount:   200,
		},
		{
			id:          5,
			name:        "PlayStation Store",
			description: "PlayStation Store gift cards add funds to your PSN wallet for games, DLC, subscriptions and more on PS4, PS5 and PS3. Global region support.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 25.00}, {Min: 11, Max: 50, Price: 24.00}, {Min: 51, Price: 23.10}},
			cardCount:   200,
		},
		{
			id:          6,
			name:        "Amazon Global",
			description: "Amazon gift cards can be used to purchase millions of items on Amazon worldwide. Valid for online purchases across all product categories.",
			rules:       []WholesaleRule{{Min: 1, Max: 10, Price: 52.00}, {Min: 11, Max: 50, Price: 50.00}, {Min: 51, Price: 48.25}},
			cardCount:   200,
		},
	}

	// 幂等：已有商品则跳过（但仍需确保测试商品存在）
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count > 0 {
		seedTestProduct(db)
		return
	}

	for _, s := range seeds {
		rulesJSON, _ := json.Marshal(s.rules)
		lowestPrice := s.rules[len(s.rules)-1].Price

		p := model.Product{
			ID:             s.id,
			CategoryID:     1,
			Name:           s.name,
			Description:    s.description,
			Price:          fmt.Sprintf("%.2f", lowestPrice),
			WholesaleRules: string(rulesJSON),
			Stock:          s.cardCount,
			Status:         1,
		}
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)

		cards := make([]model.CardKey, 0, s.cardCount)
		for i := 0; i < s.cardCount; i++ {
			cards = append(cards, model.CardKey{
				ProductID: s.id,
				Content:   generateCardContent(),
				Status:    0,
			})
		}
		db.CreateInBatches(cards, 200)
	}
}

// seedTestProduct 插入 id=99 的 0.01 测试商品（幂等）
func seedTestProduct(db *gorm.DB) {
	const testID uint64 = 99
	var exist model.Product
	if db.First(&exist, testID).Error == nil {
		return
	}
	rules := []WholesaleRule{{Min: 1, Price: 0.01}}
	rulesJSON, _ := json.Marshal(rules)
	p := model.Product{
		ID:             testID,
		CategoryID:     1,
		Name:           "Test Product ¥0.01",
		Description:    "Payment integration test item. Do not purchase in production.",
		Price:          "0.01",
		WholesaleRules: string(rulesJSON),
		Stock:          100,
		Status:         1,
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&p)
	cards := make([]model.CardKey, 100)
	for i := range cards {
		cards[i] = model.CardKey{ProductID: testID, Content: generateCardContent(), Status: 0}
	}
	db.CreateInBatches(cards, 100)
}

// ListProducts 返回所有上架产品列表
func (s *DBStore) ListProducts() []ProductDTO {
	var products []model.Product
	s.db.Where("status = 1").Order("id asc").Find(&products)

	result := make([]ProductDTO, 0, len(products))
	for _, p := range products {
		var stock int64
		s.db.Model(&model.CardKey{}).Where("product_id = ? AND status = 0", p.ID).Count(&stock)

		var rules []WholesaleRule
		_ = json.Unmarshal([]byte(p.WholesaleRules), &rules)

		result = append(result, ProductDTO{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			Price:          p.Price,
			WholesaleRules: rules,
			Stock:          int(stock),
			Status:         p.Status,
		})
	}
	return result
}

// GetProduct 根据ID返回产品
func (s *DBStore) GetProduct(id uint64) (*ProductDTO, bool) {
	var p model.Product
	if err := s.db.Where("id = ? AND status = 1", id).First(&p).Error; err != nil {
		return nil, false
	}

	var stock int64
	s.db.Model(&model.CardKey{}).Where("product_id = ? AND status = 0", id).Count(&stock)

	var rules []WholesaleRule
	_ = json.Unmarshal([]byte(p.WholesaleRules), &rules)

	return &ProductDTO{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Price:          p.Price,
		WholesaleRules: rules,
		Stock:          int(stock),
		Status:         p.Status,
	}, true
}

// CreateOrder 创建订单并锁定卡密（事务）
func (s *DBStore) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
	var out *CreateOrderOutput

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var p model.Product
		if err := tx.Where("id = ? AND status = 1", input.ProductID).First(&p).Error; err != nil {
			return fmt.Errorf("product not found")
		}

		// 锁定 N 张可用卡密（SKIP LOCKED 防止并发冲突）
		var cards []model.CardKey
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("product_id = ? AND status = 0", input.ProductID).
			Limit(input.Quantity).
			Find(&cards).Error; err != nil {
			return err
		}
		if len(cards) < input.Quantity {
			return fmt.Errorf("insufficient stock: need %d, have %d", input.Quantity, len(cards))
		}

		var rules []WholesaleRule
		_ = json.Unmarshal([]byte(p.WholesaleRules), &rules)
		unitPrice := CalcPrice(rules, input.Quantity)
		totalAmount := unitPrice * float64(input.Quantity)

		now := time.Now()
		ts := now.UnixMilli() % 100000000
		orderNo := fmt.Sprintf("SL%08d%04d", ts, rng.Intn(10000))
		expiresAt := now.Add(30 * time.Minute)

		userID := input.UserID
		if userID == 0 {
			userID = 1 // 匿名用户兜底
		}
		order := model.Order{
			OrderNo:     orderNo,
			UserID:      userID,
			ProductID:   input.ProductID,
			Quantity:    input.Quantity,
			UnitPrice:   fmt.Sprintf("%.2f", unitPrice),
			TotalAmount: fmt.Sprintf("%.2f", totalAmount),
			Status:      0,
			PayChannel:  input.PayChannel,
			Email:       input.Email,
			ExpiresAt:   expiresAt,
			CreatedAt:   now,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// 锁定卡密
		lockedAt := now
		cardIDs := make([]uint64, len(cards))
		for i, c := range cards {
			cardIDs[i] = c.ID
		}
		if err := tx.Model(&model.CardKey{}).
			Where("id IN ?", cardIDs).
			Updates(map[string]interface{}{
				"status":    1,
				"order_id":  order.ID,
				"locked_at": lockedAt,
			}).Error; err != nil {
			return err
		}

		// 写入商品快照（防止商品信息变更影响历史订单）
		snapshot := model.ProductSnapshot{
			OrderID:        order.ID,
			ProductID:      p.ID,
			Name:           p.Name,
			Price:          fmt.Sprintf("%.2f", unitPrice),
			WholesaleRules: p.WholesaleRules,
		}
		tx.Create(&snapshot)

		out = &CreateOrderOutput{
			OrderNo:     orderNo,
			TotalAmount: fmt.Sprintf("%.2f", totalAmount),
			ExpiresAt:   expiresAt,
		}
		return nil
	})

	return out, err
}

// FulfillOrder 完成支付（事务，幂等）
func (s *DBStore) FulfillOrder(orderNo string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return fmt.Errorf("order not found: %s", orderNo)
		}
		if order.Status != 0 {
			return nil // 幂等
		}
		if time.Now().After(order.ExpiresAt) {
			return fmt.Errorf("order expired")
		}

		now := time.Now()
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":  2,
			"paid_at": now,
		}).Error; err != nil {
			return err
		}

		// 查找锁定的卡密
		var cards []model.CardKey
		if err := tx.Where("order_id = ? AND status = 1", order.ID).Find(&cards).Error; err != nil {
			return err
		}

		// 标记已售，生成 order_items
		items := make([]model.OrderItem, 0, len(cards))
		cardIDs := make([]uint64, len(cards))
		for i, c := range cards {
			cardIDs[i] = c.ID
			items = append(items, model.OrderItem{
				OrderID:         order.ID,
				CardKeyID:       c.ID,
				ContentSnapshot: c.Content,
				DeliveredAt:     now,
			})
		}
		if err := tx.Model(&model.CardKey{}).Where("id IN ?", cardIDs).
			Updates(map[string]interface{}{"status": 2, "sold_at": now}).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetOrderByNo 根据订单号返回原始订单（用于支付流程）
func (s *DBStore) GetOrderByNo(orderNo string) (*model.Order, error) {
	var order model.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found: %s", orderNo)
	}
	return &order, nil
}

// GetOrderCards 获取已支付订单的卡密列表
func (s *DBStore) GetOrderCards(orderNo string) (*OrderCardsDTO, error) {
	var order model.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found: %s", orderNo)
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
