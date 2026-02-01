package services

import (
	"errors"
	"kasir-go/models"
	"kasir-go/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
	productRepo  *repositories.ProductRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository, productRepo *repositories.ProductRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo, productRepo: productRepo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *CategoryService) Create(data *models.Category) error {
	return s.categoryRepo.Create(data)
}

func (s *CategoryService) GetById(id int) (*models.Category, error) {
	return s.categoryRepo.FindById(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	products, err := s.productRepo.FindByCategoryId(id)
	if err != nil {
		return err
	}

	if len(products) > 0 {
		return errors.New("Cannot delete category: category is still used by products")
	}

	return s.categoryRepo.Delete(id)
}
