package menu

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/menu/application"
	"multicliente-backend/internal/features/menu/domain"
	"multicliente-backend/internal/features/menu/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.MenuRepository {
	repo := infrastructure.NewMenuRepository(db)
	service := application.NewMenuService(repo)
	handler := infrastructure.NewMenuHandler(service)

	// User menu query (needs to be registered outside the admin-restricted sub-group, but inside the jwt-protected router)
	router.GET("/menus/my", handler.GetMyMenus)

	menus := router.Group("/menus")
	{
		// SuperAdmin only CRUD
		menus.GET("", requireRole, handler.GetAll)
		menus.GET("/:id", requireRole, handler.GetByID)
		menus.POST("", requireRole, handler.Create)
		menus.PUT("/:id", requireRole, handler.Update)
		menus.DELETE("/:id", requireRole, handler.Delete)
	}

	return repo
}
