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

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.productRepo.FindAll(name)
}

func (s *ProductService) Create(data *models.Product) error {
	_, err := s.categoryRepo.FindById(data.CategoryID)
	if err != nil {
		return err
	}

	return s.productRepo.Create(data)
}

func (s *ProductService) GetById(id int) (*models.Product, error) {
	product, err := s.productRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.FindById(product.CategoryID)
	if err != nil {
		return nil, err
	}

	result := &models.Product{
		ID:           product.ID,
		Name:         product.Name,
		Price:        product.Price,
		CategoryID:   category.ID,
		CategoryName: category.Name,
	}

	return result, nil
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
