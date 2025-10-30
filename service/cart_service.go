package service

import (
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/model"
	"e-commerce-platform-backend/repository"
)

type CartService interface {
	GetUserCart(userId uint) (*model.Cart, error)
	AddTOCart(userID uint, cartID uint, productID uint, quantity int) error
}

type CartServiceImpl struct {
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
}

func NewDatasetServiceImpl(CartRepository repository.CartRepository, ProductRepository repository.ProductRepository) (service CartService, err error) {
	return &CartServiceImpl{
		CartRepository:    CartRepository,
		ProductRepository: ProductRepository,
	}, err
}

func (s *CartServiceImpl) GetUserCart(userId uint) (*model.Cart, error) {
	return s.CartRepository.GetUserCart(userId)
}

// GetUserCart Fetch cart with the details of the product\
func (s *CartServiceImpl) AddTOCart(userID uint, cartID uint, productID uint, quantity int) error {
	product, err := s.ProductRepository.GetProductById(productID)

	if err != nil {
		logger.ActError("Product not found")
	}

	if product.ProductStock < quantity {
		logger.ActError("Not enough stock")
	}

	//Calculate the total price of the product with the quantity
	itemPrice := product.ProductPrice * float64(quantity)

	//Adding product+price+quantity to the cart
	item := &model.CartItem{
		CartID:     cartID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  product.ProductPrice,
		TotalPrice: itemPrice,
	}

	//Add to repository
	err = s.CartRepository.AddItemToCart(item)
	if err != nil {
		logger.ActError("Error adding item to the cart")
		return err
	}

	logger.ActInfo("Items added to the cart successfully")
	return nil
}

func (s *CartServiceImpl) ClearCart(userId uint) error {
	return s.CartRepository.ClearCart(userId)
}
