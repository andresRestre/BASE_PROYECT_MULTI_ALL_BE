package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/inventory/article/application"
	"multicliente-backend/internal/features/inventory/article/domain"
	"multicliente-backend/internal/features/inventory/article/infrastructure"
	notificationDomain "multicliente-backend/internal/features/notification/domain"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc, notificationService notificationDomain.NotificationService) domain.ArticleService {
	repo := infrastructure.NewArticleRepository(db)
	service := application.NewArticleService(repo, notificationService)
	handler := infrastructure.NewArticleHandler(service)

	articles := router.Group("/articles")
	articles.Use(requireCompanyAccess)
	{
		articles.POST("", middleware.RequirePermission(db, "/inventory/items", "CREATE"), handler.Create)
		articles.GET("", middleware.RequirePermission(db, "/inventory/items", "VIEW"), handler.GetAll)
		articles.GET("/:id", middleware.RequirePermission(db, "/inventory/items", "VIEW"), handler.GetByID)
		articles.PUT("/:id", middleware.RequirePermission(db, "/inventory/items", "EDIT"), handler.Update)
		articles.DELETE("/:id", middleware.RequirePermission(db, "/inventory/items", "DELETE"), handler.Delete)
	}

	return service
}
