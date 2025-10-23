package repository

import (
	"e-commerce-platform-backend/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrderById(orderId uint) (*model.Order, error)
	UpdateOrderStatus(orderId uint, OrderStatus string) error
}

type OrderRepositoryImpl struct {
	Db *gorm.DB
}

func NewOrderRepository(Db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{Db: Db}
}

func (r OrderRepositoryImpl) GetOrderById(orderId uint) (*model.Order, error) {
	var order model.Order
	err := r.Db.Preload("Items.Products").Preload("Payment").First(&order, orderId).Error
	return &order, err
}

func (r OrderRepositoryImpl) UpdateOrderStatus(orderId uint, OrderStatus string) error {
	return r.Db.Model(&model.Order{}).
		Where("order_id=?", orderId).
		Update("order_status", OrderStatus).Error
}
