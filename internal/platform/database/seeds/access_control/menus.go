package access_control

import (
	"log"
	"gorm.io/gorm"
	menuDomain "multicliente-backend/internal/features/menu/domain"
)

// SeedMenus seeds navigation menus into database.
func SeedMenus(db *gorm.DB) error {
	accessControlID := uint(8)
	inventariosID := uint(5)
	statisticsID := uint(9)
	menus := []menuDomain.Menu{
		{ID: 8, Label: "Control de Acceso", LabelEN: "Access Control", LabelFR: "Contrôle d'Accès", Route: "/access-control", Icon: "security_rounded", SortOrder: 1, IsActive: true, ParentID: nil},
		{ID: 1, Label: "Usuarios", LabelEN: "Users", LabelFR: "Utilisateurs", Route: "/users", Icon: "people_rounded", SortOrder: 1, IsActive: true, ParentID: &accessControlID},
		{ID: 2, Label: "Empresas", LabelEN: "Companies", LabelFR: "Entreprises", Route: "/companies", Icon: "business_rounded", SortOrder: 2, IsActive: true, ParentID: &accessControlID},
		{ID: 3, Label: "Roles", LabelEN: "Roles", LabelFR: "Rôles", Route: "/roles", Icon: "admin_panel_settings_rounded", SortOrder: 3, IsActive: true, ParentID: &accessControlID},
		{ID: 4, Label: "Menús", LabelEN: "Menus", LabelFR: "Menus", Route: "/menus", Icon: "menu_rounded", SortOrder: 4, IsActive: true, ParentID: &accessControlID},
		{ID: 5, Label: "Inventarios", LabelEN: "Inventory", LabelFR: "Inventaires", Route: "/inventory", Icon: "inventory_2_rounded", SortOrder: 10, IsActive: true, ParentID: nil},
		{ID: 6, Label: "Categorías", LabelEN: "Categories", LabelFR: "Catégories", Route: "/categories", Icon: "category_rounded", SortOrder: 11, IsActive: true, ParentID: &inventariosID},
		{ID: 7, Label: "Artículos", LabelEN: "Items", LabelFR: "Articles", Route: "/items", Icon: "inventory_rounded", SortOrder: 12, IsActive: true, ParentID: &inventariosID},
		{ID: 9, Label: "Estadísticas", LabelEN: "Statistics", LabelFR: "Statistiques", Route: "/statistics", Icon: "bar_chart_rounded", SortOrder: 20, IsActive: true, ParentID: nil},
		{ID: 10, Label: "Dashboard", LabelEN: "Dashboard", LabelFR: "Tableau de Bord", Route: "/statistics/dashboard", Icon: "dashboard_rounded", SortOrder: 21, IsActive: true, ParentID: &statisticsID},
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
