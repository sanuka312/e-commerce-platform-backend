package repository

import (
	"shophub-backend/logger"
	"shophub-backend/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddItemToCart(item *model.CartItem) error
	GetUserCart(keycloakUserID string) (*model.Cart, error)
	GetOrCreateCart(keycloakUserID string) (*model.Cart, error)
	RemoveItemFromCart(itemId uint) error
	ClearCart(keycloakUserID string) error
	GetCartItemById(itemId uint) (*model.CartItem, error)
	GetCartItemByProductId(cartID uint, productID uint) (*model.CartItem, error)
	UpdateCartItemQuantity(itemId uint, quantity int) error
}

type CartRepositoryImpl struct {
	Db *gorm.DB
}

func NewCartRepository(Db *gorm.DB) CartRepository {
	return &CartRepositoryImpl{Db: Db}
}

func (r *CartRepositoryImpl) AddItemToCart(item *model.CartItem) error {
	logger.ActInfo("Adding new items to cart")
	return r.Db.Create(item).Error
}

func (r *CartRepositoryImpl) GetUserCart(keycloakUserID string) (*model.Cart, error) {
	var cart model.Cart
	query := r.Db.Model(&model.Cart{}).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_id ASC")
		}).
		Where("keycloak_user_id=?", keycloakUserID)

	if err := query.First(&cart).Error; err != nil {
		logger.ActError("Error retrieving the cart")
		return nil, err
	}

	// Ensuring products are correctly loaded and matched to cart items
	for i := range cart.Items {
		if cart.Items[i].ProductID != 0 {
			// Reload product to ensure correct matching
			var product model.Product
			if err := r.Db.First(&product, cart.Items[i].ProductID).Error; err == nil {
				cart.Items[i].Product = product
			}
		}
	}

	logger.ActInfo("Successfully retrieved user cart")
	return &cart, nil

}

func (r *CartRepositoryImpl) GetOrCreateCart(keycloakUserID string) (*model.Cart, error) {
	var cart model.Cart
	query := r.Db.Model(&model.Cart{}).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_id ASC")
		}).
		Where("keycloak_user_id=?", keycloakUserID)

	err := query.First(&cart).Error
	if err != nil {
		// Only create a new cart if the error is "record not found"
		if err == gorm.ErrRecordNotFound {
			// Cart doesn't exist, create a new one
			cart = model.Cart{
				KeycloakUserID: keycloakUserID,
				Items:          []model.CartItem{},
			}
			if err := r.Db.Create(&cart).Error; err != nil {
				logger.ActError("Error creating the cart")
				return nil, err
			}
			logger.ActInfo("Successfully created new cart")
			return &cart, nil
		}
		// For any other error, return it
		return nil, err
	}

	// Ensure products are correctly loaded and matched to cart items
	for i := range cart.Items {
		if cart.Items[i].ProductID != 0 {
			// Reload product to ensure correct matching
			var product model.Product
			if err := r.Db.First(&product, cart.Items[i].ProductID).Error; err == nil {
				cart.Items[i].Product = product
			}
		}
	}

	logger.ActInfo("Successfully retrieved user cart")
	return &cart, nil
}

func (r *CartRepositoryImpl) RemoveItemFromCart(itemId uint) error {
	return r.Db.Delete(&model.CartItem{}, itemId).Error
}

func (r *CartRepositoryImpl) ClearCart(keycloakUserID string) error {
	var cart model.Cart
	if err := r.Db.Where("keycloak_user_id=?", keycloakUserID).First(&cart).Error; err != nil {
		return err
	}

	return r.Db.Where("cart_id=?", cart.CartID).Delete(&model.CartItem{}).Error
}

func (r *CartRepositoryImpl) GetCartItemById(itemId uint) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.Db.First(&item, itemId).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepositoryImpl) GetCartItemByProductId(cartID uint, productID uint) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.Db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepositoryImpl) UpdateCartItemQuantity(itemId uint, quantity int) error {
	var item model.CartItem
	if err := r.Db.First(&item, itemId).Error; err != nil {
		return err
	}
	item.Quantity = quantity
	item.TotalPrice = item.UnitPrice * float64(quantity)
	return r.Db.Save(&item).Error
}
