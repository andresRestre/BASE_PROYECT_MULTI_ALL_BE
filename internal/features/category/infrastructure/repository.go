package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/category/domain"
	"multicliente-backend/internal/platform/database"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) populateAudits(categories []*domain.Category) {
	var userIDs []uint
	for _, c := range categories {
		if c.CreateBy != nil {
			userIDs = append(userIDs, *c.CreateBy)
		}
		if c.UpdateBy != nil {
			userIDs = append(userIDs, *c.UpdateBy)
		}
	}
	namesMap, err := database.GetUserNamesMap(r.db, userIDs)
	if err != nil {
		return
	}
	for _, c := range categories {
		if c.CreateBy != nil {
			c.CreateByName = namesMap[*c.CreateBy]
		}
		if c.UpdateBy != nil {
			c.UpdateByName = namesMap[*c.UpdateBy]
		}
	}
}

func (r *categoryRepository) Create(cat *domain.Category) error {
	return r.db.Create(cat).Error
}

func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var cat domain.Category
	if err := r.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	r.populateAudits([]*domain.Category{&cat})
	return &cat, nil
}

func (r *categoryRepository) FindAllByCompany(companyID uint) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Where("company_id = ?", companyID).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	categoriesPtrs := make([]*domain.Category, len(categories))
	for i := range categories {
		categoriesPtrs[i] = &categories[i]
	}
	r.populateAudits(categoriesPtrs)
	return categories, nil
}

func (r *categoryRepository) Update(cat *domain.Category) error {
	return r.db.Save(cat).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, "id = ?", id).Error
}
