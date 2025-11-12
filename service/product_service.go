package service

import (
	"shophub-backend/model"
	"shophub-backend/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	GetProductById(productId uint) (*model.Product, error)
}

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
}

func NewProductServiceImpl(ProductRepository repository.ProductRepository) (service ProductService, err error) {
	return &ProductServiceImpl{
		ProductRepository: ProductRepository,
	}, err
}

func (s *ProductServiceImpl) GetAllProducts() ([]model.Product, error) {
	return s.ProductRepository.GetAllProducts()
}

func (s *ProductServiceImpl) GetProductById(productId uint) (*model.Product, error) {
	return s.ProductRepository.GetProductById(productId)
}

func (s *ProductServiceImpl) GetProductBySlug(productSlug string) (*model.Product, error) {
	return s.ProductRepository.GetProductBySlug(productSlug)
}
