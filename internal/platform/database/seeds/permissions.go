package seeds

import (
	"log"
	"gorm.io/gorm"
	roleDomain "multicliente-backend/internal/features/role/domain"
)

// SeedPermissions seeds default CRUD permissions for superadmin.
func SeedPermissions(db *gorm.DB) error {
	permissions := []roleDomain.Permission{
		{RoleID: 1, MenuID: 1, OptionID: 1},
		{RoleID: 1, MenuID: 1, OptionID: 2},
		{RoleID: 1, MenuID: 1, OptionID: 3},
		{RoleID: 1, MenuID: 1, OptionID: 4},
		{RoleID: 1, MenuID: 2, OptionID: 1},
		{RoleID: 1, MenuID: 2, OptionID: 2},
		{RoleID: 1, MenuID: 2, OptionID: 3},
		{RoleID: 1, MenuID: 2, OptionID: 4},
		{RoleID: 1, MenuID: 3, OptionID: 1},
		{RoleID: 1, MenuID: 3, OptionID: 2},
		{RoleID: 1, MenuID: 3, OptionID: 3},
		{RoleID: 1, MenuID: 3, OptionID: 4},
		{RoleID: 1, MenuID: 4, OptionID: 1},
		{RoleID: 1, MenuID: 4, OptionID: 2},
		{RoleID: 1, MenuID: 4, OptionID: 3},
		{RoleID: 1, MenuID: 4, OptionID: 4},
		{RoleID: 1, MenuID: 5, OptionID: 1},
		{RoleID: 1, MenuID: 5, OptionID: 2},
		{RoleID: 1, MenuID: 5, OptionID: 3},
		{RoleID: 1, MenuID: 5, OptionID: 4},
		{RoleID: 1, MenuID: 6, OptionID: 1},
		{RoleID: 1, MenuID: 6, OptionID: 2},
		{RoleID: 1, MenuID: 6, OptionID: 3},
		{RoleID: 1, MenuID: 6, OptionID: 4},
		{RoleID: 1, MenuID: 7, OptionID: 1},
		{RoleID: 1, MenuID: 7, OptionID: 2},
		{RoleID: 1, MenuID: 7, OptionID: 3},
		{RoleID: 1, MenuID: 7, OptionID: 4},
	}
	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, roleDomain.Permission{RoleID: p.RoleID, MenuID: p.MenuID, OptionID: p.OptionID}).Error; err != nil {
			return err
		}
	}
	log.Println("✅ Permissions seeded")
	return nil
}
