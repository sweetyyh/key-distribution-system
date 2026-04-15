package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"key-distribution-system/internal/model"
)

// WholesaleRule 阶梯定价规则
type WholesaleRule struct {
	Min   int     `json:"min"`
	Max   int     `json:"max,omitempty"` // 0 表示无上限
	Price float64 `json:"price"`
}

// ProductDTO 产品数据传输对象（含解析后的规则）
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

// MemStore 内存存储
type MemStore struct {
	mu             sync.RWMutex
	products       map[uint64]*model.Product
	productRules   map[uint64][]WholesaleRule // productID -> parsed rules
	orders         map[string]*model.Order    // key: OrderNo
	cards          map[uint64]*model.CardKey
	items          map[uint64][]*model.OrderItem // key: OrderID
	cardsByProduct map[uint64][]uint64           // productID -> []available cardKeyID (FIFO)
	nextCardID     uint64
	nextOrderID    uint64
	orderEmails    map[string]string // orderNo -> email
}

// Global 全局内存存储实例
var Global = &MemStore{
	products:       make(map[uint64]*model.Product),
	productRules:   make(map[uint64][]WholesaleRule),
	orders:         make(map[string]*model.Order),
	cards:          make(map[uint64]*model.CardKey),
	items:          make(map[uint64][]*model.OrderItem),
	cardsByProduct: make(map[uint64][]uint64),
	orderEmails:    make(map[string]string),
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Init 初始化 mock 数据
func Init() {
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
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 5.00},
				{Min: 11, Max: 50, Price: 4.90},
				{Min: 51, Price: 4.75},
			},
			cardCount: 200,
		},
		{
			id:          2,
			name:        "Google Play",
			description: "Google Play gift cards let you add funds to your Google Play balance for apps, games, movies, books, and more. Multi-region support.",
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 10.00},
				{Min: 11, Max: 50, Price: 9.60},
				{Min: 51, Price: 9.20},
			},
			cardCount: 200,
		},
		{
			id:          3,
			name:        "Steam Gift Card",
			description: "Steam Gift Cards work just like a gift certificate, which can be redeemed on Steam for the purchase of games, software, wallet credit, and any other item you can purchase on Steam.",
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 50.00},
				{Min: 11, Max: 50, Price: 48.00},
				{Min: 51, Price: 45.00},
			},
			cardCount: 500,
		},
		{
			id:          4,
			name:        "Netflix Prepaid",
			description: "Netflix prepaid gift cards allow users to subscribe or extend their Netflix membership without a credit card. Supports EMEA and US regions.",
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 16.00},
				{Min: 11, Max: 50, Price: 15.50},
				{Min: 51, Price: 14.99},
			},
			cardCount: 200,
		},
		{
			id:          5,
			name:        "PlayStation Store",
			description: "PlayStation Store gift cards add funds to your PSN wallet for games, DLC, subscriptions and more on PS4, PS5 and PS3. Global region support.",
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 25.00},
				{Min: 11, Max: 50, Price: 24.00},
				{Min: 51, Price: 23.10},
			},
			cardCount: 200,
		},
		{
			id:          6,
			name:        "Amazon Global",
			description: "Amazon gift cards can be used to purchase millions of items on Amazon worldwide. Valid for online purchases across all product categories.",
			rules: []WholesaleRule{
				{Min: 1, Max: 10, Price: 52.00},
				{Min: 11, Max: 50, Price: 50.00},
				{Min: 51, Price: 48.25},
			},
			cardCount: 200,
		},
	}

	Global.mu.Lock()
	defer Global.mu.Unlock()

	for _, s := range seeds {
		// 序列化阶梯规则为 JSON
		rulesJSON, _ := json.Marshal(s.rules)

		// 用最低阶梯价作为展示价格
		lowestPrice := s.rules[len(s.rules)-1].Price

		p := &model.Product{
			ID:             s.id,
			CategoryID:     1,
			Name:           s.name,
			Description:    s.description,
			Price:          fmt.Sprintf("%.2f", lowestPrice),
			WholesaleRules: string(rulesJSON),
			Stock:          s.cardCount,
			Status:         1,
		}
		Global.products[s.id] = p
		Global.productRules[s.id] = s.rules

		// 预生成卡密
		for i := 0; i < s.cardCount; i++ {
			Global.nextCardID++
			cardID := Global.nextCardID
			content := generateCardContent()
			card := &model.CardKey{
				ID:        cardID,
				ProductID: s.id,
				Content:   content,
				Status:    0, // 未售
			}
			Global.cards[cardID] = card
			Global.cardsByProduct[s.id] = append(Global.cardsByProduct[s.id], cardID)
		}
	}
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

// CalcPrice 根据阶梯规则和数量计算单价
func CalcPrice(rules []WholesaleRule, qty int) float64 {
	var price float64
	for _, r := range rules {
		if qty >= r.Min && (r.Max == 0 || qty <= r.Max) {
			price = r.Price
			break
		}
	}
	if price == 0 && len(rules) > 0 {
		price = rules[len(rules)-1].Price
	}
	return price
}

// ListProducts 返回所有上架产品列表（按ID排序）
func (s *MemStore) ListProducts() []ProductDTO {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ProductDTO, 0, len(s.products))
	for i := uint64(1); i <= 6; i++ {
		p, ok := s.products[i]
		if !ok || p.Status != 1 {
			continue
		}
		stock := len(s.cardsByProduct[i])
		result = append(result, ProductDTO{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			Price:          p.Price,
			WholesaleRules: s.productRules[i],
			Stock:          stock,
			Status:         p.Status,
		})
	}
	return result
}

// GetProduct 根据ID返回产品
func (s *MemStore) GetProduct(id uint64) (*ProductDTO, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok || p.Status != 1 {
		return nil, false
	}
	stock := len(s.cardsByProduct[id])
	dto := &ProductDTO{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Price:          p.Price,
		WholesaleRules: s.productRules[id],
		Stock:          stock,
		Status:         p.Status,
	}
	return dto, true
}

// CreateOrder 创建订单并锁定卡密
func (s *MemStore) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.products[input.ProductID]
	if !ok || p.Status != 1 {
		return nil, fmt.Errorf("product not found")
	}

	available := s.cardsByProduct[input.ProductID]
	if len(available) < input.Quantity {
		return nil, fmt.Errorf("insufficient stock: need %d, have %d", input.Quantity, len(available))
	}

	rules := s.productRules[input.ProductID]
	unitPrice := CalcPrice(rules, input.Quantity)
	totalAmount := unitPrice * float64(input.Quantity)

	// 生成订单号
	s.nextOrderID++
	ts := time.Now().UnixMilli() % 100000000
	orderNo := fmt.Sprintf("SL%08d%04d", ts, s.nextOrderID%10000)

	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)

	order := &model.Order{
		ID:          s.nextOrderID,
		OrderNo:     orderNo,
		UserID:      1, // 匿名用户
		ProductID:   input.ProductID,
		Quantity:    input.Quantity,
		UnitPrice:   fmt.Sprintf("%.2f", unitPrice),
		TotalAmount: fmt.Sprintf("%.2f", totalAmount),
		Status:      0, // 待支付
		PayChannel:  input.PayChannel,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}
	s.orders[orderNo] = order
	s.orderEmails[orderNo] = input.Email

	// 锁定卡密（取前 Quantity 张）
	tolock := available[:input.Quantity]
	s.cardsByProduct[input.ProductID] = available[input.Quantity:]

	lockedAt := now
	for _, cardID := range tolock {
		card := s.cards[cardID]
		card.Status = 1 // 锁定
		card.OrderID = &order.ID
		card.LockedAt = &lockedAt
	}

	return &CreateOrderOutput{
		OrderNo:     orderNo,
		TotalAmount: fmt.Sprintf("%.2f", totalAmount),
		ExpiresAt:   expiresAt,
	}, nil
}

// FulfillOrder 完成支付，将订单和卡密标记为已完成
func (s *MemStore) FulfillOrder(orderNo string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderNo]
	if !ok {
		return fmt.Errorf("order not found: %s", orderNo)
	}
	if order.Status != 0 {
		// 幂等：已完成则直接返回
		return nil
	}
	if time.Now().After(order.ExpiresAt) {
		return fmt.Errorf("order expired")
	}

	now := time.Now()
	order.Status = 2 // 已完成
	order.PaidAt = &now

	// 找到该订单锁定的卡密，标记已售，创建 OrderItem
	var orderItems []*model.OrderItem
	itemID := uint64(0)
	for _, card := range s.cards {
		if card.OrderID != nil && *card.OrderID == order.ID && card.Status == 1 {
			soldAt := now
			card.Status = 2 // 已售
			card.SoldAt = &soldAt

			itemID++
			item := &model.OrderItem{
				ID:              itemID,
				OrderID:         order.ID,
				CardKeyID:       card.ID,
				ContentSnapshot: card.Content,
				DeliveredAt:     now,
			}
			orderItems = append(orderItems, item)
		}
	}
	s.items[order.ID] = orderItems

	return nil
}

// GetOrderCards 获取订单的卡密列表
func (s *MemStore) GetOrderCards(orderNo string) (*OrderCardsDTO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderNo]
	if !ok {
		return nil, fmt.Errorf("order not found: %s", orderNo)
	}
	if order.Status != 2 {
		return nil, fmt.Errorf("order not paid yet")
	}

	p := s.products[order.ProductID]
	productName := ""
	if p != nil {
		productName = p.Name
	}

	items := s.items[order.ID]
	cards := make([]OrderCardItem, 0, len(items))
	for i, item := range items {
		// ContentSnapshot 格式: SL-XXXX-XXXX-XXXX|DDDD
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

// splitContent 从 "code|pin" 格式中分离 code 和 pin
func splitContent(content string) (string, string) {
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '|' {
			return content[:i], content[i+1:]
		}
	}
	return content, ""
}
