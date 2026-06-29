package category

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/category/application"
	"multicliente-backend/internal/features/category/domain"
	"multicliente-backend/internal/features/category/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc) domain.CategoryService {
	repo := infrastructure.NewCategoryRepository(db)
	service := application.NewCategoryService(repo)
	handler := infrastructure.NewCategoryHandler(service)

	categories := router.Group("/categories")
	categories.Use(requireCompanyAccess)
	{
		categories.POST("", handler.Create)
		categories.GET("", handler.GetAll)
		categories.GET("/:id", handler.GetByID)
		categories.PUT("/:id", handler.Update)
		categories.DELETE("/:id", handler.Delete)
	}

	return service
}
