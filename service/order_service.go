package service

import (
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/model"
	"e-commerce-platform-backend/repository"
	"errors"
	"time"
)

type OrderService interface {
	CreateOrder(userId uint) (*model.Order, error)
}

type OrderServiceImpl struct {
	OrderRepository   repository.OrderRepository
	ProductRepository repository.ProductRepository
	CartRepository    repository.CartRepository
	PaymentRepository repository.PaymentRepository
}

func NewOrderServiceImpl(OrderRepository repository.OrderRepository, ProductRepository repository.ProductRepository, CartRepository repository.CartRepository, PaymentRepository repository.PaymentRepository) (service OrderService, err error) {
	return &OrderServiceImpl{
		OrderRepository:   OrderRepository,
		CartRepository:    CartRepository,
		ProductRepository: ProductRepository,
		PaymentRepository: PaymentRepository,
	}, err
}

func (s *OrderServiceImpl) CreateOrder(userId uint) (*model.Order, error) {
	cart, err := s.CartRepository.GetUserCart(userId)
	if err != nil || cart == nil {
		logger.ActError("Cart not found")
		return nil, err
	}

	if len(cart.Items) == 0 {
		logger.ActError("cart is empty")

	}

	order := &model.Order{
		UserId:      userId,
		TotalPrice:  0,
		OrderStatus: "Pending",
		CreatedAt:   time.Now(),
	}

	for _, item := range cart.Items {
		product, err := s.ProductRepository.GetProductById(item.ProductID)

		if err != nil {
			return nil, err
		}

		if product.ProductStock < item.Quantity {
			return nil, errors.New("Insufficient stock for " + product.ProductName)
		}

		product.ProductStock -= item.Quantity

		s.ProductRepository.UpdateProduct(product)

		order.TotalPrice += float64(item.Quantity) * product.ProductPrice
	}

	if err := s.OrderRepository.CreateOrder(order); err != nil {
		logger.ActError("Unable to create the order")
		return nil, err
	}

	if err := s.CartRepository.ClearCart(userId); err != nil {
		return nil, err
	}

	payment := &model.Payment{
		OrderId:       order.OrderId,
		UserId:        userId,
		PaymentAmount: order.TotalPrice,
		Status:        "UNPAID",
	}

	s.PaymentRepository.CreatePayment(payment)

	return order, nil
}

func (s *OrderServiceImpl) GetOrderByUser(userId uint) ([]model.Order, error) {
	return s.OrderRepository.GetOrderByUserId(userId)
}
