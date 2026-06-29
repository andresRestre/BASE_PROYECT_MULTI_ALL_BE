package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/menu/domain"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domain.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(menu *domain.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindByID(id uint) (*domain.Menu, error) {
	var menu domain.Menu
	if err := r.db.First(&menu, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll() ([]domain.Menu, error) {
	var menus []domain.Menu
	if err := r.db.Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) Update(menu *domain.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Menu{}, "id = ?", id).Error
}

func (r *menuRepository) GetAllowedMenusForRole(roleID uint) ([]domain.AllowedMenuResponse, error) {
	type resultRow struct {
		ID        uint
		Label     string
		LabelEN   string
		Route     string
		Icon      string
		SortOrder int
		ParentID  *uint
		OptCode   string
	}

	var rows []resultRow
	err := r.db.Table("administrative.menus m").
		Select("m.id, m.label, m.label_en, m.route, m.icon, m.sort_order, m.parent_id, o.code as opt_code").
		Joins("JOIN administrative.permissions p ON p.menu_id = m.id").
		Joins("JOIN administrative.options o ON o.id = p.option_id").
		Where("p.role_id = ? AND m.is_active = true", roleID).
		Order("m.sort_order ASC, m.id ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// 1. Group rows by menu ID to aggregate permissions list
	menusMap := make(map[uint]*domain.AllowedMenuResponse)
	for _, row := range rows {
		if _, ok := menusMap[row.ID]; !ok {
			menusMap[row.ID] = &domain.AllowedMenuResponse{
				ID:          row.ID,
				Label:       row.Label,
				LabelEN:     row.LabelEN,
				Route:       row.Route,
				Icon:        row.Icon,
				SortOrder:   row.SortOrder,
				ParentID:    row.ParentID,
				Permissions: []string{},
				Submenus:    []domain.AllowedMenuResponse{},
			}
		}
		menusMap[row.ID].Permissions = append(menusMap[row.ID].Permissions, row.OptCode)
	}

	// 2. Load missing parent folders from DB to ensure navigation tree remains continuous
	for _, menu := range menusMap {
		if menu.ParentID != nil {
			pID := *menu.ParentID
			if _, ok := menusMap[pID]; !ok {
				var parentMenu domain.Menu
				if err := r.db.First(&parentMenu, "id = ?", pID).Error; err == nil {
					menusMap[pID] = &domain.AllowedMenuResponse{
						ID:          parentMenu.ID,
						Label:       parentMenu.Label,
						LabelEN:     parentMenu.LabelEN,
						Route:       parentMenu.Route,
						Icon:        parentMenu.Icon,
						SortOrder:   parentMenu.SortOrder,
						ParentID:    parentMenu.ParentID,
						Permissions: []string{},
						Submenus:    []domain.AllowedMenuResponse{},
					}
				}
			}
		}
	}

	// 3. Separate root menus and child menus
	var rootMenus []*domain.AllowedMenuResponse
	for _, menu := range menusMap {
		if menu.ParentID == nil {
			rootMenus = append(rootMenus, menu)
		}
	}

	// 4. Populate Submenus
	for _, menu := range menusMap {
		if menu.ParentID != nil {
			pID := *menu.ParentID
			if parent, ok := menusMap[pID]; ok {
				parent.Submenus = append(parent.Submenus, *menu)
			}
		}
	}

	// Helper to sort AllowedMenuResponse slice in-place
	sortResponses := func(slice []domain.AllowedMenuResponse) {
		for i := 0; i < len(slice); i++ {
			for j := i + 1; j < len(slice); j++ {
				if slice[i].SortOrder > slice[j].SortOrder || 
				   (slice[i].SortOrder == slice[j].SortOrder && slice[i].ID > slice[j].ID) {
					slice[i], slice[j] = slice[j], slice[i]
				}
			}
		}
	}

	// 5. Build and sort final response list
	var finalResponses []domain.AllowedMenuResponse
	for _, root := range rootMenus {
		sortResponses(root.Submenus)
		finalResponses = append(finalResponses, *root)
	}

	sortResponses(finalResponses)

	return finalResponses, nil
}
