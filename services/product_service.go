package services

import (
	"kasir-go/models"
	"kasir-go/repositories"
)

type ProductService struct {
	productRepo  *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetAll() ([]models.ProductResponse, error) {
	return s.productRepo.FindAll()
}

func (s *ProductService) Create(data *models.Product) error {
	_, err := s.categoryRepo.FindById(data.CategoryID)
	if err != nil {
		return err
	}

	return s.productRepo.Create(data)
}

func (s *ProductService) GetById(id int) (*models.ProductDetailResponse, error) {
	return s.productRepo.FindById(id)
}

func (s *ProductService) Update(product *models.Product) error {
	_, err := s.categoryRepo.FindById(product.CategoryID)
	if err != nil {
		return err
	}

	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(id)
}
