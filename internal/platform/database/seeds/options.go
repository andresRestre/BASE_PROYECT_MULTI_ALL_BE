package seeds

import (
	"log"
	"gorm.io/gorm"
	roleDomain "multicliente-backend/internal/features/role/domain"
)

// SeedOptions seeds option types (CRUD) into database.
func SeedOptions(db *gorm.DB) error {
	options := []roleDomain.Option{
		{ID: 1, Name: "Ver", Code: "VIEW"},
		{ID: 2, Name: "Editar", Code: "EDIT"},
		{ID: 3, Name: "Crear", Code: "CREATE"},
		{ID: 4, Name: "Eliminar", Code: "DELETE"},
	}
	for _, opt := range options {
		if err := db.FirstOrCreate(&opt, roleDomain.Option{ID: opt.ID}).Error; err != nil {
			return err
		}
	}
	log.Println("✅ Options seeded")
	return nil
}
