package access_control

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	companyDomain "multicliente-backend/internal/features/company/domain"
	userDomain "multicliente-backend/internal/features/user/domain"
)

// SeedUsers seeds default administrator user.
func SeedUsers(db *gorm.DB) error {
	var count int64
	if err := db.Model(&userDomain.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		adminRoleID := uint(1)
		companies := []companyDomain.Company{
			{ID: 1, Name: "Empresa Base Demo", IsActive: true},
		}

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
			return err
		}
		log.Println("✅ Admin user seeded successfully (admin@example.com / admin123)")
	}
	return nil
}
