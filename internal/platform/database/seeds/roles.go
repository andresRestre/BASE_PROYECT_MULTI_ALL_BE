package seeds

import (
	"log"
	"gorm.io/gorm"
	roleDomain "multicliente-backend/internal/features/role/domain"
)

// SeedRoles seeds system roles into database.
func SeedRoles(db *gorm.DB) error {
	roles := []roleDomain.Role{
		{ID: 1, Name: "Super Administrador", Code: "superadmin", IsActive: true},
		{ID: 2, Name: "Administrador", Code: "admin", IsActive: true},
		{ID: 3, Name: "Usuario", Code: "user", IsActive: true},
	}
	for _, r := range roles {
		if err := db.FirstOrCreate(&r, roleDomain.Role{ID: r.ID}).Error; err != nil {
			return err
		}
	}
	log.Println("✅ Roles seeded")
	return nil
}
