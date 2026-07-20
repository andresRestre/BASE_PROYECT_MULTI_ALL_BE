package access_control

import (
	"log"
	"gorm.io/gorm"
	menuDomain "multicliente-backend/internal/features/access_control/menu/domain"
	roleDomain "multicliente-backend/internal/features/access_control/role/domain"
)

// SeedPermissions seeds default CRUD permissions for all menus dynamically for SuperAdmin and Admin roles.
func SeedPermissions(db *gorm.DB) error {
	var menus []menuDomain.Menu
	if err := db.Find(&menus).Error; err != nil {
		return err
	}

	roleIDs := []uint{1, 2} // Grant full permissions to SuperAdmin (1) and Admin (2)
	optionIDs := []uint{1, 2, 3, 4} // VIEW, CREATE, EDIT, DELETE

	for _, menu := range menus {
		for _, roleID := range roleIDs {
			for _, optionID := range optionIDs {
				perm := roleDomain.Permission{
					RoleID:   roleID,
					MenuID:   menu.ID,
					OptionID: optionID,
				}
				if err := db.FirstOrCreate(&perm, roleDomain.Permission{
					RoleID:   roleID,
					MenuID:   menu.ID,
					OptionID: optionID,
				}).Error; err != nil {
					return err
				}
			}
		}
	}

	log.Println("✅ Dynamic Menu Permissions Seeded Successfully")
	return nil
}
