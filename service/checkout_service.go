package service

import (
	"errors"
	"shophub-backend/data"
	"shophub-backend/logger"
	"shophub-backend/model"
	"shophub-backend/repository"
	"time"

	"go.uber.org/zap"
)

type CheckoutService interface {
	PlaceOrder(keycloakUserID string, paymentMethod string, address data.CreateAddressRequest) (*model.Order, error)
}

type CheckoutServiceImpl struct {
	OrderRepository   repository.OrderRepository
	ProductRepository repository.ProductRepository
	CartRepository    repository.CartRepository
	PaymentRepository repository.PaymentRepository
	AddressRepository repository.AddressRepository
	UserRepository    repository.UserRepository
}

func NewCheckoutServiceImpl(
	OrderRepository repository.OrderRepository,
	ProductRepository repository.ProductRepository,
	CartRepository repository.CartRepository,
	PaymentRepository repository.PaymentRepository,
	AddressRepository repository.AddressRepository,
	UserRepository repository.UserRepository,
) (CheckoutService, error) {
	return &CheckoutServiceImpl{
		OrderRepository:   OrderRepository,
		ProductRepository: ProductRepository,
		CartRepository:    CartRepository,
		PaymentRepository: PaymentRepository,
		AddressRepository: AddressRepository,
		UserRepository:    UserRepository,
	}, nil
}

func (s *CheckoutServiceImpl) PlaceOrder(keycloakUserID string, paymentMethod string, addressReq data.CreateAddressRequest) (*model.Order, error) {
	// Get user's cart
	cart, err := s.CartRepository.GetUserCart(keycloakUserID)
	if err != nil || cart == nil {
		logger.ActError("Cart not found")
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		logger.ActError("cart is empty")
		return nil, errors.New("cart is empty")
	}

	// Validate stock for all items before processing any orders
	for _, item := range cart.Items {
		product, err := s.ProductRepository.GetProductById(item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.ProductStock < item.Quantity {
			return nil, errors.New("Insufficient stock for " + product.ProductName)
		}
	}

	// Ensure user exists in database (required for foreign key constraint)
	// This creates a minimal user record if it doesn't exist
	_, err = s.UserRepository.GetOrCreateUser(keycloakUserID)
	if err != nil {
		logger.ActError("Unable to ensure user exists", zap.Error(err))
		return nil, errors.New("failed to ensure user exists: " + err.Error())
	}

	// Create address for the order
	address := &model.Address{
		KeycloakUserID: keycloakUserID,
		Line1:          addressReq.Line1,
		Line2:          addressReq.Line2,
		City:           addressReq.City,
		PostalCode:     addressReq.PostalCode,
		Country:        addressReq.Country,
	}

	if err := s.AddressRepository.CreateAddress(address); err != nil {
		logger.ActError("Unable to create address for order", zap.Error(err))
		return nil, errors.New("failed to create address: " + err.Error())
	}

	// Verify address ID was populated
	if address.AddressId == 0 {
		logger.ActError("Address ID not populated after creation")
		return nil, errors.New("address ID not populated after creation")
	}

	// Normalize payment method (convert to uppercase for consistency)
	normalizedPaymentMethod := paymentMethod
	if normalizedPaymentMethod == "Cash on Delivery" || normalizedPaymentMethod == "cash on delivery" {
		normalizedPaymentMethod = "CASH"
	} else if normalizedPaymentMethod == "Credit/Debit Card" || normalizedPaymentMethod == "credit/debit card" {
		normalizedPaymentMethod = "CARD"
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

		// Create payment without OrderId (will be updated after order creation)
		payment := &model.Payment{
			OrderId:        nil, // Set to nil initially, will be updated after order creation
			KeycloakUserID: keycloakUserID,
			PaymentMethod:  normalizedPaymentMethod,
			PaymentAmount:  itemTotalPrice,
			Status:         "UNPAID",
		}

		// Creating the payment for the order
		if err := s.PaymentRepository.CreatePayment(payment); err != nil {
			logger.ActError("Unable to create payment for order", zap.Error(err))
			return nil, errors.New("failed to create payment: " + err.Error())
		}

		// Create order with all required fields including address
		order := &model.Order{
			KeycloakUserID: keycloakUserID,
			ProductId:      item.ProductID,
			PaymentId:      payment.PaymentId,
			ProductPrice:   product.ProductPrice,
			Quantity:       uint(item.Quantity),
			TotalPrice:     itemTotalPrice,
			AddressId:      &address.AddressId,
			OrderStatus:    "Pending",
			CreatedAt:      time.Now(),
		}

		if err := s.OrderRepository.CreateOrder(order); err != nil {
			logger.ActError("Unable to create the order", zap.Error(err))
			return nil, errors.New("failed to create order: " + err.Error())
		}

		// Update payment with the actual OrderId
		if err := s.PaymentRepository.UpdatePaymentOrderId(payment.PaymentId, order.OrderId); err != nil {
			logger.ActError("Unable to update payment with order ID", zap.Error(err))
			return nil, errors.New("failed to update payment order ID: " + err.Error())
		}

		// IMPORTANT: Reduce stock after order is successfully created
		// This ensures stock is only reduced when the order is confirmed
		product.ProductStock -= item.Quantity
		if err := s.ProductRepository.UpdateProduct(product); err != nil {
			logger.ActError("Unable to update product stock", zap.Error(err))
			return nil, errors.New("failed to update product stock: " + err.Error())
		}

		// Store first order to return
		if userOrder == nil {
			userOrder = order
		}
	}

	// Clear cart after all orders are created
	if err := s.CartRepository.ClearCart(keycloakUserID); err != nil {
		logger.ActError("Unable to clear cart after order creation", zap.Error(err))
		return nil, errors.New("failed to clear cart: " + err.Error())
	}

	return userOrder, nil
}
