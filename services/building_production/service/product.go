package service

import (
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/building_production/repository/mariadb"
)

type ProductService interface {
	GetDefProducts() ([]model.Product, error)
	AddProduct(product model.Product) error
}

type productService struct {
	productRepo mariadb.ProductRepository
}

func NewProductService(productRepo mariadb.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetDefProducts() ([]model.Product, error) {
	return s.productRepo.GetDefProducts()
}

func (s *productService) AddProduct(product model.Product) error {
	return s.productRepo.CreateProduct(product)
}
