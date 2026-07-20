package inventory

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"multicliente-backend/internal/features/inventory/article"
	"multicliente-backend/internal/features/inventory/category"
	notificationDomain "multicliente-backend/internal/features/notification/domain"
)

// RegisterRoutes registers all inventory feature routes (category, article) under a parent route group.
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc, notificationService notificationDomain.NotificationService) {
	category.RegisterRoutes(router, db, requireCompanyAccess, notificationService)
	article.RegisterRoutes(router, db, requireCompanyAccess, notificationService)
}
