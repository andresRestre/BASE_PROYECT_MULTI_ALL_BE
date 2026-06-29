package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/company/domain"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) domain.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(company *domain.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) FindByID(id uint) (*domain.Company, error) {
	var company domain.Company
	if err := r.db.First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) FindAll() ([]domain.Company, error) {
	var companies []domain.Company
	if err := r.db.Order("create_at DESC").Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) Update(company *domain.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Company{}, "id = ?", id).Error
}
