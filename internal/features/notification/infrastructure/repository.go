package infrastructure

import (
	"strings"

	"gorm.io/gorm"
	"multicliente-backend/internal/features/notification/domain"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindAllByUserAndCompany(userID uint, companyID uint) ([]domain.Notification, error) {
	var list []domain.Notification
	err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *notificationRepository) FindByID(id uint) (*domain.Notification, error) {
	var notif domain.Notification
	if err := r.db.First(&notif, id).Error; err != nil {
		return nil, err
	}
	return &notif, nil
}

func (r *notificationRepository) Update(notification *domain.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint, companyID uint) error {
	return r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND company_id = ? AND is_read = ?", userID, companyID, false).
		Update("is_read", true).Error
}

func (r *notificationRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Notification{}).Error
}

func (r *notificationRepository) DeleteAllByUserAndCompany(userID uint, companyID uint) error {
	return r.db.Where("user_id = ? AND company_id = ?", userID, companyID).Delete(&domain.Notification{}).Error
}

func (r *notificationRepository) FindAdminsAndSuperadminsByCompany(companyID uint) ([]uint, error) {
	// 1. All superadmins (who manage all companies in system)
	var superAdminIDs []uint
	_ = r.db.Table("administrative.users u").
		Select("u.id").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("LOWER(r.code) = ? AND u.is_active = true", "superadmin").
		Pluck("u.id", &superAdminIDs).Error

	// 2. Admins associated with the specific company
	var companyAdminIDs []uint
	_ = r.db.Table("administrative.users u").
		Select("u.id").
		Joins("JOIN administrative.user_companies uc ON uc.user_id = u.id").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("uc.company_id = ? AND (LOWER(r.code) = ? OR LOWER(r.code) LIKE ?) AND u.is_active = true", companyID, "admin", "%admin%").
		Pluck("u.id", &companyAdminIDs).Error

	idMap := make(map[uint]bool)
	for _, id := range superAdminIDs {
		idMap[id] = true
	}
	for _, id := range companyAdminIDs {
		idMap[id] = true
	}

	var result []uint
	for id := range idMap {
		result = append(result, id)
	}
	return result, nil
}

func (r *notificationRepository) GetCreatorInfo(userID uint) (string, string, string, error) {
	type Result struct {
		FirstName string
		LastName  string
		RoleCode  string
	}
	var res Result
	err := r.db.Table("administrative.users u").
		Select("u.first_name, u.last_name, r.code as role_code").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("u.id = ?", userID).
		Scan(&res).Error
	if err != nil {
		return "", "", "", err
	}
	return res.FirstName, res.LastName, res.RoleCode, nil
}

func (r *notificationRepository) GetActorDetails(userID uint) (*domain.ActorDetails, error) {
	type Result struct {
		FirstName string
		LastName  string
		RoleID    uint
		RoleName  string
		RoleCode  string
	}
	var res Result
	err := r.db.Table("administrative.users u").
		Select("u.first_name, u.last_name, r.id as role_id, r.name as role_name, r.code as role_code").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("u.id = ?", userID).
		Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return &domain.ActorDetails{
		FirstName: res.FirstName,
		LastName:  res.LastName,
		RoleID:    res.RoleID,
		RoleName:  res.RoleName,
		RoleCode:  res.RoleCode,
	}, nil
}

func (r *notificationRepository) FindTargetUsersByNotificationRule(companyID uint, menuRoute string, creatorRoleID uint, action string) ([]uint, error) {
	shortRoute := menuRoute
	if idx := strings.LastIndex(menuRoute, "/"); idx >= 0 {
		shortRoute = menuRoute[idx:]
	}
	var menu struct {
		ID uint
	}
	if err := r.db.Table("administrative.menus").Select("id").
		Where("route = ? OR route = ?", menuRoute, shortRoute).
		First(&menu).Error; err != nil {
		return nil, err
	}

	var targetRoleIDs []uint
	if err := r.db.Table("administrative.role_notification_rules").
		Select("target_role_id").
		Where("menu_id = ? AND creator_role_id = ? AND UPPER(action) = ? AND is_enabled = true", menu.ID, creatorRoleID, strings.ToUpper(action)).
		Pluck("target_role_id", &targetRoleIDs).Error; err != nil {
		return nil, err
	}

	if len(targetRoleIDs) == 0 {
		return nil, nil
	}

	var userIDs []uint
	_ = r.db.Table("administrative.users u").
		Select("u.id").
		Where("u.role_id IN ? AND u.is_active = true", targetRoleIDs).
		Pluck("u.id", &userIDs).Error

	idMap := make(map[uint]bool)
	var result []uint
	for _, id := range userIDs {
		if !idMap[id] {
			idMap[id] = true
			result = append(result, id)
		}
	}
	return result, nil
}
