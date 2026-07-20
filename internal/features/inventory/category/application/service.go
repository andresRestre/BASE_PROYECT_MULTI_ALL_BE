package application

import (
	"errors"
	"multicliente-backend/internal/features/inventory/category/domain"
)

type categoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req *domain.CreateCategoryRequest, companyID uint, createdBy *uint) (*domain.Category, error) {
	cat := &domain.Category{
		CompanyID: companyID,
		Name:      req.Name,
		IsActive:  true,
		CreateBy:  createdBy,
	}

	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *categoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (s *categoryService) GetCategoriesByCompany(companyID uint) ([]domain.Category, error) {
	return s.repo.FindAllByCompany(companyID)
}

func (s *categoryService) UpdateCategory(id uint, req *domain.UpdateCategoryRequest, updatedBy *uint) (*domain.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	cat.UpdateBy = updatedBy

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(id)
}
