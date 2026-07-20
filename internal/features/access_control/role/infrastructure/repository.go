package infrastructure

import (
	"gorm.io/gorm"
	"multicliente-backend/internal/features/access_control/role/domain"
	"multicliente-backend/internal/platform/database"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) populateAudits(roles []*domain.Role) {
	var userIDs []uint
	for _, r := range roles {
		if r.CreateBy != nil {
			userIDs = append(userIDs, *r.CreateBy)
		}
		if r.UpdateBy != nil {
			userIDs = append(userIDs, *r.UpdateBy)
		}
	}
	namesMap, err := database.GetUserNamesMap(r.db, userIDs)
	if err != nil {
		return
	}
	for _, r := range roles {
		if r.CreateBy != nil {
			r.CreateByName = namesMap[*r.CreateBy]
		}
		if r.UpdateBy != nil {
			r.UpdateByName = namesMap[*r.UpdateBy]
		}
	}
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Preload("Permissions").Preload("NotificationRules").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	r.populateAudits([]*domain.Role{&role})
	return &role, nil
}

func (r *roleRepository) FindAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.db.Preload("Permissions").Preload("NotificationRules").Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	rolesPtrs := make([]*domain.Role, len(roles))
	for i := range roles {
		rolesPtrs[i] = &roles[i]
	}
	r.populateAudits(rolesPtrs)
	return roles, nil
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Role{}, "id = ?", id).Error
}

func (r *roleRepository) ReplacePermissions(roleID uint, permissions []domain.Permission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing permissions
		if err := tx.Delete(&domain.Permission{}, "role_id = ?", roleID).Error; err != nil {
			return err
		}
		// Insert new permissions if any
		if len(permissions) > 0 {
			if err := tx.Create(&permissions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) ReplaceNotificationRules(roleID uint, rules []domain.RoleNotificationRule) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing notification rules for this target role
		if err := tx.Delete(&domain.RoleNotificationRule{}, "target_role_id = ?", roleID).Error; err != nil {
			return err
		}
		// Insert new notification rules
		if len(rules) > 0 {
			for i := range rules {
				rules[i].TargetRoleID = roleID
			}
			if err := tx.Create(&rules).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) FindAllOptions() ([]domain.Option, error) {
	var options []domain.Option
	if err := r.db.Order("id ASC").Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *roleRepository) GetUserRoleHierarchy(userID uint) (int, error) {
	var hierarchy int
	err := r.db.Table("administrative.users u").
		Select("r.hierarchy").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("u.id = ?", userID).
		Scan(&hierarchy).Error
	if err != nil {
		return 999, err
	}
	return hierarchy, nil
}
