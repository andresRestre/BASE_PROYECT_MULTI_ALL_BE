package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/category/application"
	"multicliente-backend/internal/features/category/domain"
	"multicliente-backend/internal/features/category/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc) domain.CategoryService {
	repo := infrastructure.NewCategoryRepository(db)
	service := application.NewCategoryService(repo)
	handler := infrastructure.NewCategoryHandler(service)

	categories := router.Group("/categories")
	categories.Use(requireCompanyAccess)
	{
		categories.POST("", middleware.RequirePermission(db, "/categories", "CREATE"), handler.Create)
		categories.GET("", middleware.RequirePermission(db, "/categories", "VIEW"), handler.GetAll)
		categories.GET("/:id", middleware.RequirePermission(db, "/categories", "VIEW"), handler.GetByID)
		categories.PUT("/:id", middleware.RequirePermission(db, "/categories", "EDIT"), handler.Update)
		categories.DELETE("/:id", middleware.RequirePermission(db, "/categories", "DELETE"), handler.Delete)
	}

	return service
}
