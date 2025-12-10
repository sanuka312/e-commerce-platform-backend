package service

import (
	"fmt"
	"shophub-backend/logger"
	"shophub-backend/model"
	"shophub-backend/repository"
)

type CartService interface {
	GetUserCart(keycloakUserID string) (*model.Cart, error)
	AddTOCart(keycloakUserID string, productID uint, quantity int) error
	ClearCart(keycloakUserID string) error
	RemoveItemFromCart(itemId uint) error
	UpdateCartItemQuantity(itemId uint, quantity int) error
}

type CartServiceImpl struct {
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
}

func NewCartServiceImpl(CartRepository repository.CartRepository, ProductRepository repository.ProductRepository) (service CartService, err error) {
	return &CartServiceImpl{
		CartRepository:    CartRepository,
		ProductRepository: ProductRepository,
	}, err
}

func (s *CartServiceImpl) GetUserCart(keycloakUserID string) (*model.Cart, error) {
	// Use GetOrCreateCart to ensure a cart always exists (even if empty)
	return s.CartRepository.GetOrCreateCart(keycloakUserID)
}

// AddTOCart adds items to cart with the details of the product
func (s *CartServiceImpl) AddTOCart(keycloakUserID string, productID uint, quantity int) error {
	// Get or create cart for the user
	cart, err := s.CartRepository.GetOrCreateCart(keycloakUserID)
	if err != nil {
		logger.ActError("Error getting or creating cart")
		return fmt.Errorf("failed to get or create cart")
	}

	product, err := s.ProductRepository.GetProductById(productID)
	if err != nil {
		logger.ActError("Product not found")
		return fmt.Errorf("product not found")
	}

	// Check if item already exists in cart
	existingItem, err := s.CartRepository.GetCartItemByProductId(cart.CartID, productID)
	if err == nil {
		// Item exists, update quantity instead of creating duplicate
		newQuantity := existingItem.Quantity + quantity

		// Validate stock availability for the new total quantity
		if product.ProductStock < newQuantity {
			logger.ActError("Not enough stock")
			return fmt.Errorf("insufficient stock for the product. Only %d item(s) available", product.ProductStock)
		}

		// Update the existing item's quantity
		err = s.CartRepository.UpdateCartItemQuantity(existingItem.ID, newQuantity)
		if err != nil {
			logger.ActError("Error updating cart item quantity")
			return fmt.Errorf("failed to update cart item quantity")
		}

		logger.ActInfo("Cart item quantity updated successfully")
		return nil
	}

	// Validate stock availability
	if product.ProductStock < quantity {
		logger.ActError("Not enough stock")
		return fmt.Errorf("insufficient stock for the product")
	}

	//Calculate the total price of the product with the quantity
	itemPrice := product.ProductPrice * float64(quantity)

	//Adding product+price+quantity to the cart
	item := &model.CartItem{
		CartID:     cart.CartID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  product.ProductPrice,
		TotalPrice: itemPrice,
	}

	//Adding item to the cart
	err = s.CartRepository.AddItemToCart(item)
	if err != nil {
		logger.ActError("Error adding item to the cart")
		return err
	}

	logger.ActInfo("Items added to the cart successfully")
	return nil
}

func (s *CartServiceImpl) ClearCart(keycloakUserID string) error {
	return s.CartRepository.ClearCart(keycloakUserID)
}

func (s *CartServiceImpl) RemoveItemFromCart(itemId uint) error {
	return s.CartRepository.RemoveItemFromCart(itemId)

}

func (s *CartServiceImpl) UpdateCartItemQuantity(itemId uint, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	// Get the cart item to find the product ID
	cartItem, err := s.CartRepository.GetCartItemById(itemId)
	if err != nil {
		logger.ActError("Cart item not found")
		return fmt.Errorf("cart item not found")
	}

	// Get the product to check stock availability
	product, err := s.ProductRepository.GetProductById(cartItem.ProductID)
	if err != nil {
		logger.ActError("Product not found")
		return fmt.Errorf("product not found")
	}

	// Validate stock availability
	if product.ProductStock < quantity {
		logger.ActError("Not enough stock")
		return fmt.Errorf("insufficient stock for the product. Only %d item(s) available", product.ProductStock)
	}

	// Update the quantity
	return s.CartRepository.UpdateCartItemQuantity(itemId, quantity)
}
