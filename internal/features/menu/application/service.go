package application

import (
	"errors"
	"multicliente-backend/internal/features/menu/domain"
)

type menuService struct {
	repo domain.MenuRepository
}

func NewMenuService(repo domain.MenuRepository) domain.MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) CreateMenu(req *domain.CreateMenuRequest, createdBy *uint) (*domain.Menu, error) {
	menu := &domain.Menu{
		Label:     req.Label,
		LabelEN:   req.LabelEN,
		Route:     req.Route,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		ParentID:  req.ParentID,
		IsActive:  true,
		CreateBy:  createdBy,
	}

	if menu.ParentID != nil && *menu.ParentID == 0 {
		menu.ParentID = nil
	}

	if err := s.repo.Create(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *menuService) GetMenuByID(id uint) (*domain.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	return menu, nil
}

func (s *menuService) GetAllMenus() ([]domain.Menu, error) {
	return s.repo.FindAll()
}

func (s *menuService) UpdateMenu(id uint, req *domain.UpdateMenuRequest, updatedBy *uint) (*domain.Menu, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	if req.Label != nil {
		menu.Label = *req.Label
	}
	if req.LabelEN != nil {
		menu.LabelEN = *req.LabelEN
	}
	if req.Route != nil {
		menu.Route = *req.Route
	}
	if req.Icon != nil {
		menu.Icon = *req.Icon
	}
	if req.IsActive != nil {
		menu.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		menu.SortOrder = *req.SortOrder
	}
	if req.ParentID != nil {
		if *req.ParentID == 0 {
			menu.ParentID = nil
		} else {
			menu.ParentID = req.ParentID
		}
	}
	menu.UpdateBy = updatedBy

	if err := s.repo.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *menuService) DeleteMenu(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("menu not found")
	}
	return s.repo.Delete(id)
}

func (s *menuService) GetAllowedMenus(roleID uint) ([]domain.AllowedMenuResponse, error) {
	return s.repo.GetAllowedMenusForRole(roleID)
}
