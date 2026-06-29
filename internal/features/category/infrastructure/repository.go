package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/category/domain"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(cat *domain.Category) error {
	return r.db.Create(cat).Error
}

func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var cat domain.Category
	if err := r.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindAllByCompany(companyID uint) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Where("company_id = ?", companyID).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Update(cat *domain.Category) error {
	return r.db.Save(cat).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, "id = ?", id).Error
}
