package service

import (
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/model"
	"e-commerce-platform-backend/repository"
)

type PaymentService interface {
	GetPaymentByOrderId(OrderId uint) (*model.Payment, error)
	ProcessPayment(orderId uint, paymentMethod string) (*model.Payment, error)
}

type PaymentServiceImpl struct {
	PaymentRepository repository.PaymentRepository
	OrderRepository   repository.OrderRepository
}

func NewPaymentServiceImpl(PaymentRepository repository.PaymentRepository, OrderRepository repository.OrderRepository) (service PaymentService, err error) {
	return &PaymentServiceImpl{
		PaymentRepository: PaymentRepository,
		OrderRepository:   OrderRepository,
	}, err
}

func (s *PaymentServiceImpl) GetPaymentByOrderId(OrderId uint) (*model.Payment, error) {
	return s.PaymentRepository.GetPaymentByOrder(OrderId)
}

func (s *PaymentServiceImpl) ProcessPayment(orderId uint, paymentMethod string) (*model.Payment, error) {
	order, err := s.OrderRepository.GetOrderById(orderId)

	if err != nil {
		logger.ActError("Unable to find the order")
	}

	payment, err := s.PaymentRepository.GetPaymentByOrder(orderId)

	if err != nil {
		logger.ActError("error occured while fetching the payment")
	}

	payment.Status = "PAID"
	payment.PaymentMethod = "paymentMethod"

	if err := s.PaymentRepository.UpdatePaymentStatus(payment.OrderId, "PAID"); err != nil {
		return nil, err
	}

	s.OrderRepository.UpdateOrderStatus(order.OrderId, "CONFIRMED")

	return payment, nil

}
