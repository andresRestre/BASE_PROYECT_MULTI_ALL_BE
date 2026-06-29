package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	companyDomain "multicliente-backend/internal/features/company/domain"
	menuDomain "multicliente-backend/internal/features/menu/domain"
	roleDomain "multicliente-backend/internal/features/role/domain"
	userDomain "multicliente-backend/internal/features/user/domain"
)

// SeedDatabase populates default options, roles, menus, permissions, companies, and the admin user.
func SeedDatabase(db *gorm.DB) {
	// 1. Seed Options
	options := []roleDomain.Option{
		{ID: 1, Name: "Ver", Code: "VIEW"},
		{ID: 2, Name: "Editar", Code: "EDIT"},
	}
	for _, opt := range options {
		db.FirstOrCreate(&opt, roleDomain.Option{ID: opt.ID})
	}
	log.Println("✅ Options seeded")

	// 2. Seed Roles
	roles := []roleDomain.Role{
		{ID: 1, Name: "Super Administrador", Code: "superadmin", IsActive: true},
		{ID: 2, Name: "Administrador", Code: "admin", IsActive: true},
		{ID: 3, Name: "Usuario", Code: "user", IsActive: true},
	}
	for _, r := range roles {
		db.FirstOrCreate(&r, roleDomain.Role{ID: r.ID})
	}
	log.Println("✅ Roles seeded")

	// 3. Seed Menus
	inventariosID := uint(5)
	menus := []menuDomain.Menu{
		{ID: 1, Label: "Usuarios", LabelEN: "Users", Route: "/users", Icon: "people_rounded", SortOrder: 1, IsActive: true},
		{ID: 2, Label: "Empresas", LabelEN: "Companies", Route: "/companies", Icon: "business_rounded", SortOrder: 2, IsActive: true},
		{ID: 3, Label: "Roles", LabelEN: "Roles", Route: "/roles", Icon: "admin_panel_settings_rounded", SortOrder: 3, IsActive: true},
		{ID: 4, Label: "Menús", LabelEN: "Menus", Route: "/menus", Icon: "menu_rounded", SortOrder: 4, IsActive: true},
		{ID: 5, Label: "Inventarios", LabelEN: "Inventory", Route: "/inventory", Icon: "inventory_2_rounded", SortOrder: 10, IsActive: true, ParentID: nil},
		{ID: 6, Label: "Categorías", LabelEN: "Categories", Route: "/categories", Icon: "category_rounded", SortOrder: 11, IsActive: true, ParentID: &inventariosID},
		{ID: 7, Label: "Artículos", LabelEN: "Items", Route: "/items", Icon: "inventory_rounded", SortOrder: 12, IsActive: true, ParentID: &inventariosID},
	}
	for _, m := range menus {
		db.FirstOrCreate(&m, menuDomain.Menu{ID: m.ID})
	}
	log.Println("✅ Menus seeded")

	// 4. Seed Permissions (superadmin gets VIEW and EDIT for all menus)
	permissions := []roleDomain.Permission{
		{RoleID: 1, MenuID: 1, OptionID: 1},
		{RoleID: 1, MenuID: 1, OptionID: 2},
		{RoleID: 1, MenuID: 2, OptionID: 1},
		{RoleID: 1, MenuID: 2, OptionID: 2},
		{RoleID: 1, MenuID: 3, OptionID: 1},
		{RoleID: 1, MenuID: 3, OptionID: 2},
		{RoleID: 1, MenuID: 4, OptionID: 1},
		{RoleID: 1, MenuID: 4, OptionID: 2},
		{RoleID: 1, MenuID: 5, OptionID: 1},
		{RoleID: 1, MenuID: 5, OptionID: 2},
		{RoleID: 1, MenuID: 6, OptionID: 1},
		{RoleID: 1, MenuID: 6, OptionID: 2},
		{RoleID: 1, MenuID: 7, OptionID: 1},
		{RoleID: 1, MenuID: 7, OptionID: 2},
	}
	for _, p := range permissions {
		db.FirstOrCreate(&p, roleDomain.Permission{RoleID: p.RoleID, MenuID: p.MenuID, OptionID: p.OptionID})
	}
	log.Println("✅ Permissions seeded")

	// 5. Seed Default Company
	companies := []companyDomain.Company{
		{ID: 1, Name: "Empresa Base Demo", IsActive: true},
	}
	for _, c := range companies {
		db.FirstOrCreate(&c, companyDomain.Company{ID: c.ID})
	}
	log.Println("✅ Default company seeded")

	// 6. Seed Default Admin User if empty
	var count int64
	db.Model(&userDomain.User{}).Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash seed password: %v", err)
			return
		}

		adminRoleID := uint(1)
		admin := &userDomain.User{
			Email:     "admin@example.com",
			Password:  string(hashedPassword),
			FirstName: "Admin",
			LastName:  "User",
			IsActive:  true,
			RoleID:    &adminRoleID,
			Companies: companies,
		}

		if err := db.Create(admin).Error; err != nil {
			log.Printf("Failed to seed admin user: %v", err)
			return
		}
		log.Println("✅ Admin user seeded successfully (admin@example.com / admin123)")
	}
}
