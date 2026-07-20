package infrastructure

import (
	"multicliente-backend/internal/features/notification/domain"
	"gorm.io/gorm"
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

func (r *notificationRepository) FindAdminsAndSuperadminsByCompany(companyID uint) ([]uint, error) {
	var userIDs []uint
	err := r.db.Table("administrative.users u").
		Select("u.id").
		Joins("JOIN administrative.user_companies uc ON uc.user_id = u.id").
		Joins("JOIN administrative.roles r ON r.id = u.role_id").
		Where("uc.company_id = ? AND (r.code = ? OR r.code = ?)", companyID, "admin", "superadmin").
		Pluck("u.id", &userIDs).Error
	return userIDs, err
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
