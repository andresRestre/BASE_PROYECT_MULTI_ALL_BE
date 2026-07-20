package application

import (
	"errors"
	"multicliente-backend/internal/features/access_control/role/domain"
)

type roleService struct {
	repo domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) domain.RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(req *domain.CreateRoleRequest, createdBy *uint) (*domain.Role, error) {
	if createdBy != nil {
		actorHierarchy, _ := s.repo.GetUserRoleHierarchy(*createdBy)
		if req.Hierarchy < actorHierarchy {
			return nil, errors.New("no tienes permiso para crear un rol con una jerarquía de mayor poder a la tuya")
		}
	}

	role := &domain.Role{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		Hierarchy:      req.Hierarchy,
		SessionDays:    req.SessionDays,
		SessionHours:   req.SessionHours,
		SessionMinutes: req.SessionMinutes,
		IsActive:       true,
		CreateBy:       createdBy,
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

	// Add notification rules if provided
	if len(req.NotificationRules) > 0 {
		rules := make([]domain.RoleNotificationRule, len(req.NotificationRules))
		for i, r := range req.NotificationRules {
			rules[i] = domain.RoleNotificationRule{
				TargetRoleID:  role.ID,
				MenuID:        r.MenuID,
				CreatorRoleID: r.CreatorRoleID,
				Action:        r.Action,
				IsEnabled:     r.IsEnabled,
			}
		}
		if err := s.repo.ReplaceNotificationRules(role.ID, rules); err != nil {
			return nil, err
		}
		role.NotificationRules = rules
	}

	return s.repo.FindByID(role.ID)
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

	if updatedBy != nil {
		actorHierarchy, _ := s.repo.GetUserRoleHierarchy(*updatedBy)
		if role.Hierarchy < actorHierarchy {
			return nil, errors.New("no tienes permiso para modificar un rol con mayor jerarquía a la tuya")
		}
		if req.Hierarchy != nil && *req.Hierarchy < actorHierarchy {
			return nil, errors.New("no tienes permiso para asignar a un rol una jerarquía de mayor poder a la tuya")
		}
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Code != nil {
		role.Code = *req.Code
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Hierarchy != nil {
		role.Hierarchy = *req.Hierarchy
	}
	if req.SessionDays != nil {
		role.SessionDays = *req.SessionDays
	}
	if req.SessionHours != nil {
		role.SessionHours = *req.SessionHours
	}
	if req.SessionMinutes != nil {
		role.SessionMinutes = *req.SessionMinutes
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

	// Always sync/replace notification rules during update
	rules := make([]domain.RoleNotificationRule, len(req.NotificationRules))
	for i, r := range req.NotificationRules {
		rules[i] = domain.RoleNotificationRule{
			TargetRoleID:  role.ID,
			MenuID:        r.MenuID,
			CreatorRoleID: r.CreatorRoleID,
			Action:        r.Action,
			IsEnabled:     r.IsEnabled,
		}
	}
	if err := s.repo.ReplaceNotificationRules(role.ID, rules); err != nil {
		return nil, err
	}

	return s.repo.FindByID(role.ID)
}

func (s *roleService) DeleteRole(id uint, deletedBy *uint) error {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	if deletedBy != nil {
		actorHierarchy, _ := s.repo.GetUserRoleHierarchy(*deletedBy)
		if role.Hierarchy < actorHierarchy {
			return errors.New("no tienes permiso para eliminar un rol con mayor jerarquía a la tuya")
		}
	}

	return s.repo.Delete(id)
}

func (s *roleService) GetAllOptions() ([]domain.Option, error) {
	return s.repo.FindAllOptions()
}
