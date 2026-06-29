package application

import (
	"errors"
	"multicliente-backend/internal/features/role/domain"
)

type roleService struct {
	repo domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) domain.RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(req *domain.CreateRoleRequest, createdBy *uint) (*domain.Role, error) {
	role := &domain.Role{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: true,
		CreateBy: createdBy,
	}

	if err := s.repo.Create(role); err != nil {
		return nil, err
	}

	// Add permissions if provided
	if len(req.Permissions) > 0 {
		perms := make([]domain.Permission, len(req.Permissions))
		for i, p := range req.Permissions {
			perms[i] = domain.Permission{
				RoleID:   role.ID,
				MenuID:   p.MenuID,
				OptionID: p.OptionID,
			}
		}
		if err := s.repo.ReplacePermissions(role.ID, perms); err != nil {
			return nil, err
		}
		role.Permissions = perms
	}

	return role, nil
}

func (s *roleService) GetRoleByID(id uint) (*domain.Role, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (s *roleService) GetAllRoles() ([]domain.Role, error) {
	return s.repo.FindAll()
}

func (s *roleService) UpdateRole(id uint, req *domain.UpdateRoleRequest, updatedBy *uint) (*domain.Role, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Code != nil {
		role.Code = *req.Code
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}
	role.UpdateBy = updatedBy

	if err := s.repo.Update(role); err != nil {
		return nil, err
	}

	// Always sync/replace permissions during update
	perms := make([]domain.Permission, len(req.Permissions))
	for i, p := range req.Permissions {
		perms[i] = domain.Permission{
			RoleID:   role.ID,
			MenuID:   p.MenuID,
			OptionID: p.OptionID,
		}
	}
	if err := s.repo.ReplacePermissions(role.ID, perms); err != nil {
		return nil, err
	}
	role.Permissions = perms

	return role, nil
}

func (s *roleService) DeleteRole(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}
	return s.repo.Delete(id)
}

func (s *roleService) GetAllOptions() ([]domain.Option, error) {
	return s.repo.FindAllOptions()
}
