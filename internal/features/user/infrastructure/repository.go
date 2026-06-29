package infrastructure

import (
	"gorm.io/gorm"

	"multicliente-backend/internal/features/user/domain"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-based UserRepository.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Role").Preload("Companies").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Role").Preload("Companies").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload("Role").Preload("Companies").Order("create_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Save basic user details
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		// Sync many-to-many Companies association
		// By doing this within transaction, we ensure the join table matches user.Companies
		if err := tx.Model(user).Association("Companies").Replace(user.Companies); err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user domain.User
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			return err
		}

		// Clear company associations first
		if err := tx.Model(&user).Association("Companies").Clear(); err != nil {
			return err
		}

		// Delete the user
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
}
