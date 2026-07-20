package application

import (
	"fmt"
	"multicliente-backend/internal/features/notification/domain"
)

type notificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) domain.NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(companyID uint, userID uint, title, message, notifType, route string) (*domain.Notification, error) {
	notif := &domain.Notification{
		CompanyID: companyID,
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      notifType,
		Route:     route,
		IsRead:    false,
	}
	if err := s.repo.Create(notif); err != nil {
		return nil, err
	}
	return notif, nil
}

func (s *notificationService) GetNotifications(userID uint, companyID uint) ([]domain.Notification, error) {
	return s.repo.FindAllByUserAndCompany(userID, companyID)
}

func (s *notificationService) MarkAsRead(id uint, userID uint) (*domain.Notification, error) {
	notif, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if notif.UserID != userID {
		return nil, fmt.Errorf("unauthorized to update this notification")
	}
	notif.IsRead = true
	if err := s.repo.Update(notif); err != nil {
		return nil, err
	}
	return notif, nil
}

func (s *notificationService) MarkAllRead(userID uint, companyID uint) error {
	return s.repo.MarkAllAsRead(userID, companyID)
}

func (s *notificationService) TriggerArticleCreatedNotification(companyID uint, creatorID uint, articleName string) error {
	firstName, lastName, roleCode, err := s.repo.GetCreatorInfo(creatorID)
	if err != nil {
		return err
	}

	// Trigger notifications only if the creator is NOT an admin or superadmin
	if roleCode == "admin" || roleCode == "superadmin" {
		return nil
	}

	// Fetch all admins and superadmins in that company
	adminIDs, err := s.repo.FindAdminsAndSuperadminsByCompany(companyID)
	if err != nil {
		return err
	}

	title := "Nuevo Artículo Creado"
	message := fmt.Sprintf("El usuario %s %s ha creado el artículo: %s", firstName, lastName, articleName)

	for _, adminID := range adminIDs {
		// Avoid notifying oneself
		if adminID == creatorID {
			continue
		}
		_, _ = s.CreateNotification(companyID, adminID, title, message, "article_created", "/inventory/items")
	}

	return nil
}
