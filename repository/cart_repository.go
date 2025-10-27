package repository

import (
	"e-commerce-platform-backend/logger"
	"e-commerce-platform-backend/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddItemToCart(item *model.CartItem) error
	GetUserCart(userId uint) (*model.Cart, error)
	RemoveItemFromCart(itemId uint) error
	ClearCart(userId uint) error
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

func (r *CartRepositoryImpl) GetUserCart(userId uint) (*model.Cart, error) {
	logger.ActInfo("Getting user cart")
	var cart model.Cart
	query := r.Db.Model(&model.Cart{}).
		Preload("Items.Product").
		Where("user_id=?", userId)

	if err := query.First(&cart).Error; err != nil {
		logger.ActError("Error retrieving the cart")
		return nil, err
	}
	logger.ActInfo("Successfully retrieved user cart")
	return &cart, nil

}

func (r *CartRepositoryImpl) RemoveItemFromCart(itemId uint) error {
	return r.Db.Delete(&model.CartItem{}, itemId).Error
}

func (r *CartRepositoryImpl) ClearCart(userId uint) error {
	var cart model.Cart
	if err := r.Db.Where("user_id=?", userId).First(&cart).Error; err != nil {
		return err
	}

	return r.Db.Where("cart_id=?", cart.CartID).Delete(&model.CartItem{}).Error
}
