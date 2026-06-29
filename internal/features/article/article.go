package article

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/article/application"
	"multicliente-backend/internal/features/article/domain"
	"multicliente-backend/internal/features/article/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireCompanyAccess gin.HandlerFunc) domain.ArticleService {
	repo := infrastructure.NewArticleRepository(db)
	service := application.NewArticleService(repo)
	handler := infrastructure.NewArticleHandler(service)

	articles := router.Group("/articles")
	articles.Use(requireCompanyAccess)
	{
		articles.POST("", handler.Create)
		articles.GET("", handler.GetAll)
		articles.GET("/:id", handler.GetByID)
		articles.PUT("/:id", handler.Update)
		articles.DELETE("/:id", handler.Delete)
	}

	return service
}
