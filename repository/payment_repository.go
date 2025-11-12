package repository

import (
	"shophub-backend/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *model.Payment) error
	GetPaymentByOrder(orderId uint) (*model.Payment, error)
	UpdatePaymentStatus(orderId uint, status string) error
	UpdatePaymentOrderId(paymentId uint, orderId uint) error
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

// updating the payment status
func (r *PaymentRepositoryImpl) UpdatePaymentStatus(orderId uint, status string) error {
	return r.Db.Model(&model.Payment{}).
		Where("order_id=?", orderId).
		Update("payment_status", status).Error
}

func (r *PaymentRepositoryImpl) UpdatePaymentOrderId(paymentId uint, orderId uint) error {
	return r.Db.Model(&model.Payment{}).
		Where("payment_id=?", paymentId).
		Update("order_id", orderId).Error
}
