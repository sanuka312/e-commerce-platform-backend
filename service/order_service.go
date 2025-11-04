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
	GetOrderByUser(userId uint) ([]model.Order, error)
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
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		logger.ActError("cart is empty")
		return nil, errors.New("cart is empty")
	}

	// validating stock for all items before processing any orders
	for _, item := range cart.Items {
		product, err := s.ProductRepository.GetProductById(item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.ProductStock < item.Quantity {
			return nil, errors.New("Insufficient stock for " + product.ProductName)
		}
	}

	// Create one order per product
	var userOrder *model.Order
	for _, item := range cart.Items {
		product, err := s.ProductRepository.GetProductById(item.ProductID)
		if err != nil {
			return nil, err
		}

		// Calculate price for this item
		itemTotalPrice := float64(item.Quantity) * product.ProductPrice

		// Create payment with temporary OrderId (will be updated after order creation)
		payment := &model.Payment{
			OrderId:       0,
			UserId:        userId,
			PaymentMethod: "CASH",
			PaymentAmount: itemTotalPrice,
			Status:        "UNPAID",
		}

		//Creating the payment for the order
		if err := s.PaymentRepository.CreatePayment(payment); err != nil {
			logger.ActError("Unable to create payment for order")
			return nil, err
		}

		// Create order with all required fields
		order := &model.Order{
			UserId:       userId,
			ProductId:    item.ProductID,
			PaymentId:    payment.PaymentId,
			ProductPrice: product.ProductPrice,
			Quantity:     uint(item.Quantity),
			TotalPrice:   itemTotalPrice,
			OrderStatus:  "Pending",
			CreatedAt:    time.Now(),
		}

		if err := s.OrderRepository.CreateOrder(order); err != nil {
			logger.ActError("Unable to create the order")
			return nil, err
		}

		// Update payment with the actual OrderId
		if err := s.PaymentRepository.UpdatePaymentOrderId(payment.PaymentId, order.OrderId); err != nil {
			logger.ActError("Unable to update payment with order ID")
			return nil, err
		}

		// Update stock after order is created
		product.ProductStock -= item.Quantity
		if err := s.ProductRepository.UpdateProduct(product); err != nil {
			logger.ActError("Unable to update product stock")
			return nil, err
		}

		// Store first order to return
		if userOrder == nil {
			userOrder = order
		}
	}

	// Clear cart after all orders are created
	if err := s.CartRepository.ClearCart(userId); err != nil {
		logger.ActError("Unable to clear cart after order creation")
		return nil, err
	}

	return userOrder, nil
}

func (s *OrderServiceImpl) GetOrderByUser(userId uint) ([]model.Order, error) {
	return s.OrderRepository.GetOrderByUserId(userId)
}
