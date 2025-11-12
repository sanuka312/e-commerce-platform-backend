package repository

import (
	"shophub-backend/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	GetOrderById(orderId uint) (*model.Order, error)
	GetOrderByUserId(UserId uint) ([]model.Order, error)
	UpdateOrderStatus(orderId uint, OrderStatus string) error
}

type OrderRepositoryImpl struct {
	Db *gorm.DB
}

func NewOrderRepository(Db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{Db: Db}
}

func (r OrderRepositoryImpl) CreateOrder(order *model.Order) error {
	return r.Db.Create(order).Error
}

func (r OrderRepositoryImpl) GetOrderByUserId(UserId uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.Db.
		Preload("Product").
		Preload("Payment").
		Where("user_id=?", UserId).
		Find(&orders).Error
	return orders, err
}

// Getting order with the payment and products
func (r OrderRepositoryImpl) GetOrderById(orderId uint) (*model.Order, error) {
	var order model.Order
	err := r.Db.Preload("Product").Preload("Payment").First(&order, orderId).Error
	return &order, err
}

func (r OrderRepositoryImpl) UpdateOrderStatus(orderId uint, OrderStatus string) error {
	return r.Db.Model(&model.Order{}).
		Where("order_id=?", orderId).
		Update("order_status", OrderStatus).Error
}
