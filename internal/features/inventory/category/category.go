package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/inventory/category/application"
	"multicliente-backend/internal/features/inventory/category/domain"
	"multicliente-backend/internal/features/inventory/category/infrastructure"
	notificationDomain "multicliente-backend/internal/features/notification/domain"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc, notificationService notificationDomain.NotificationService) domain.CategoryService {
	repo := infrastructure.NewCategoryRepository(db)
	service := application.NewCategoryService(repo, notificationService)
	handler := infrastructure.NewCategoryHandler(service)

	categories := router.Group("/categories")
	categories.Use(requireCompanyAccess)
	{
		categories.POST("", middleware.RequirePermission(db, "/inventory/categories", "CREATE"), handler.Create)
		categories.GET("", middleware.RequirePermission(db, "/inventory/categories", "VIEW"), handler.GetAll)
		categories.GET("/:id", middleware.RequirePermission(db, "/inventory/categories", "VIEW"), handler.GetByID)
		categories.PUT("/:id", middleware.RequirePermission(db, "/inventory/categories", "EDIT"), handler.Update)
		categories.DELETE("/:id", middleware.RequirePermission(db, "/inventory/categories", "DELETE"), handler.Delete)
	}

	return service
}
