package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/article/application"
	"multicliente-backend/internal/features/article/domain"
	"multicliente-backend/internal/features/article/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc) domain.ArticleService {
	repo := infrastructure.NewArticleRepository(db)
	service := application.NewArticleService(repo)
	handler := infrastructure.NewArticleHandler(service)

	articles := router.Group("/articles")
	articles.Use(requireCompanyAccess)
	{
		articles.POST("", middleware.RequirePermission(db, "/items", "CREATE"), handler.Create)
		articles.GET("", middleware.RequirePermission(db, "/items", "VIEW"), handler.GetAll)
		articles.GET("/:id", middleware.RequirePermission(db, "/items", "VIEW"), handler.GetByID)
		articles.PUT("/:id", middleware.RequirePermission(db, "/items", "EDIT"), handler.Update)
		articles.DELETE("/:id", middleware.RequirePermission(db, "/items", "DELETE"), handler.Delete)
	}

	return service
}
