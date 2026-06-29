package seeds

import (
	"log"
	"gorm.io/gorm"
	menuDomain "multicliente-backend/internal/features/menu/domain"
)

// SeedMenus seeds navigation menus into database.
func SeedMenus(db *gorm.DB) error {
	inventariosID := uint(5)
	menus := []menuDomain.Menu{
		{ID: 1, Label: "Usuarios", LabelEN: "Users", LabelFR: "Utilisateurs", Route: "/users", Icon: "people_rounded", SortOrder: 1, IsActive: true},
		{ID: 2, Label: "Empresas", LabelEN: "Companies", LabelFR: "Entreprises", Route: "/companies", Icon: "business_rounded", SortOrder: 2, IsActive: true},
		{ID: 3, Label: "Roles", LabelEN: "Roles", LabelFR: "Rôles", Route: "/roles", Icon: "admin_panel_settings_rounded", SortOrder: 3, IsActive: true},
		{ID: 4, Label: "Menús", LabelEN: "Menus", LabelFR: "Menus", Route: "/menus", Icon: "menu_rounded", SortOrder: 4, IsActive: true},
		{ID: 5, Label: "Inventarios", LabelEN: "Inventory", LabelFR: "Inventaires", Route: "/inventory", Icon: "inventory_2_rounded", SortOrder: 10, IsActive: true, ParentID: nil},
		{ID: 6, Label: "Categorías", LabelEN: "Categories", LabelFR: "Catégories", Route: "/categories", Icon: "category_rounded", SortOrder: 11, IsActive: true, ParentID: &inventariosID},
		{ID: 7, Label: "Artículos", LabelEN: "Items", LabelFR: "Articles", Route: "/items", Icon: "inventory_rounded", SortOrder: 12, IsActive: true, ParentID: &inventariosID},
	}
	for _, m := range menus {
		var existing menuDomain.Menu
		if err := db.Where("id = ?", m.ID).First(&existing).Error; err != nil {
			if err := db.Create(&m).Error; err != nil {
				return err
			}
		} else {
			existing.Label = m.Label
			existing.LabelEN = m.LabelEN
			existing.LabelFR = m.LabelFR
			existing.Route = m.Route
			existing.Icon = m.Icon
			existing.SortOrder = m.SortOrder
			existing.ParentID = m.ParentID
			if err := db.Save(&existing).Error; err != nil {
				return err
			}
		}
	}
	log.Println("✅ Menus seeded")
	return nil
}
