package application

import (
	"errors"
	"multicliente-backend/internal/features/company/domain"
)

type companyService struct {
	repo domain.CompanyRepository
}

func NewCompanyService(repo domain.CompanyRepository) domain.CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(req *domain.CreateCompanyRequest, createdBy *uint) (*domain.Company, error) {
	company := &domain.Company{
		Name:     req.Name,
		IsActive: true,
		CreateBy: createdBy,
		PhotoURL: req.PhotoURL,
	}

	if err := s.repo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyService) GetCompanyByID(id uint) (*domain.Company, error) {
	company, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("company not found")
	}
	return company, nil
}

func (s *companyService) GetAllCompanies() ([]domain.Company, error) {
	return s.repo.FindAll()
}

func (s *companyService) UpdateCompany(id uint, req *domain.UpdateCompanyRequest, updatedBy *uint) (*domain.Company, error) {
	company, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("company not found")
	}

	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}
	if req.PhotoURL != nil {
		company.PhotoURL = *req.PhotoURL
	}
	company.UpdateBy = updatedBy

	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyService) DeleteCompany(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("company not found")
	}
	return s.repo.Delete(id)
}
