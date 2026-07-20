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
	return s.TriggerEntityEventNotification(companyID, creatorID, "/inventory/items", "CREATE", articleName)
}

func (s *notificationService) TriggerEntityEventNotification(companyID uint, actorID uint, menuRoute string, action string, entityName string) error {
	actor, err := s.repo.GetActorDetails(actorID)
	actorName := "Usuario"
	actorRoleName := "Usuario"
	actorRoleID := uint(0)
	if err == nil && actor != nil {
		actorName = fmt.Sprintf("%s %s", actor.FirstName, actor.LastName)
		actorRoleName = actor.RoleName
		actorRoleID = actor.RoleID
	}

	// Find targeted user IDs via dynamic RoleNotificationRule matrix
	targetUserIDs, _ := s.repo.FindTargetUsersByNotificationRule(companyID, menuRoute, actorRoleID, action)

	// Fallback: if no specific rules exist yet, default to notifying admins and superadmins
	if len(targetUserIDs) == 0 {
		targetUserIDs, _ = s.repo.FindAdminsAndSuperadminsByCompany(companyID)
	}

	actionVerb := "procesado"
	actionTitle := "Registro Actualizado"
	switch action {
	case "CREATE":
		actionVerb = "creado"
		actionTitle = "Nuevo Registro Creado"
	case "EDIT":
		actionVerb = "modificado"
		actionTitle = "Registro Modificado"
	case "DELETE":
		actionVerb = "eliminado"
		actionTitle = "Registro Eliminado"
	}

	title := actionTitle
	message := fmt.Sprintf("El usuario %s (%s) ha %s: %s", actorName, actorRoleName, actionVerb, entityName)

	for _, targetID := range targetUserIDs {
		if targetID == actorID {
			continue // Avoid notifying oneself
		}
		_, _ = s.CreateNotification(companyID, targetID, title, message, "entity_event", menuRoute)
	}

	return nil
}
