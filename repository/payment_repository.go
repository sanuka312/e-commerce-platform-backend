package repository

import (
	"e-commerce-platform-backend/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	GetPaymentByOrder(orderId uint) (*model.Payment, error)
	UpdatePaymentStatus(orderId uint, status string) error
}

type PaymentRepositoryImpl struct {
	Db *gorm.DB
}

func NewPaymentRepositoryImpl(Db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{Db: Db}
}

func (r *PaymentRepositoryImpl) CreatePayment(payment *model.Payment) error {
	return r.Db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) GetPaymentByOrder(orderId uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.Db.Where("order_id=?", orderId).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepositoryImpl) UpdatePaymentStatus(orderId uint, status string) error {
	return r.Db.Model(&model.Payment{}).
		Where("order_id=?", orderId).
		Update("payment_status", status).Error
}
