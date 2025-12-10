package repository

import (
	"shophub-backend/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *model.Product) error
	GetAllProducts() ([]model.Product, error)
	GetProductById(productId uint) (*model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(productID uint) error
	GetProductBySlug(productSlug string) (*model.Product, error)
}

type ProductRepositoryImpl struct {
	Db *gorm.DB
}

func NewProductRepository(Db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{Db: Db}
}

func (r ProductRepositoryImpl) CreateProduct(product *model.Product) error {
	return r.Db.Create(product).Error
}

func (r ProductRepositoryImpl) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := r.Db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r ProductRepositoryImpl) GetProductById(productId uint) (*model.Product, error) {
	var product model.Product
	if err := r.Db.Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
		return db.Order("image_id ASC")
	}).First(&product, productId).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r ProductRepositoryImpl) UpdateProduct(product *model.Product) error {
	return r.Db.Save(product).Error
}

func (r ProductRepositoryImpl) DeleteProduct(productID uint) error {
	return r.Db.Delete(&model.Product{}, productID).Error
}

func (r ProductRepositoryImpl) GetProductBySlug(productSlug string) (*model.Product, error) {
	var product model.Product
	if err := r.Db.Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
		return db.Order("image_id ASC")
	}).
		Where("product_slug= ?", productSlug).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
