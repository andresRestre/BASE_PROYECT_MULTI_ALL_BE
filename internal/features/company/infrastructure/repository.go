package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/company/domain"
	"multicliente-backend/internal/platform/database"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) domain.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) populateAudits(companies []*domain.Company) {
	var userIDs []uint
	for _, c := range companies {
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
	for _, c := range companies {
		if c.CreateBy != nil {
			c.CreateByName = namesMap[*c.CreateBy]
		}
		if c.UpdateBy != nil {
			c.UpdateByName = namesMap[*c.UpdateBy]
		}
	}
}

func (r *companyRepository) Create(company *domain.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) FindByID(id uint) (*domain.Company, error) {
	var company domain.Company
	if err := r.db.First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	r.populateAudits([]*domain.Company{&company})
	return &company, nil
}

func (r *companyRepository) FindAll() ([]domain.Company, error) {
	var companies []domain.Company
	if err := r.db.Order("create_at DESC").Find(&companies).Error; err != nil {
		return nil, err
	}
	companiesPtrs := make([]*domain.Company, len(companies))
	for i := range companies {
		companiesPtrs[i] = &companies[i]
	}
	r.populateAudits(companiesPtrs)
	return companies, nil
}

func (r *companyRepository) Update(company *domain.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Company{}, "id = ?", id).Error
}
