package notification

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/notification/application"
	"multicliente-backend/internal/features/notification/domain"
	"multicliente-backend/internal/features/notification/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) domain.NotificationService {
	repo := infrastructure.NewNotificationRepository(db)
	service := application.NewNotificationService(repo)
	handler := infrastructure.NewNotificationHandler(service)

	notifs := router.Group("/notifications")
	{
		notifs.GET("", handler.GetNotifications)
		notifs.PUT("/:id/read", handler.MarkAsRead)
		notifs.PUT("/mark-all-read", handler.MarkAllAsRead)
		notifs.DELETE("/:id", handler.DeleteNotification)
		notifs.DELETE("", handler.DeleteAllNotifications)
	}

	return service
}
