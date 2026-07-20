package menu

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/access_control/menu/application"
	"multicliente-backend/internal/features/access_control/menu/domain"
	"multicliente-backend/internal/features/access_control/menu/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.MenuRepository {
	repo := infrastructure.NewMenuRepository(db)
	service := application.NewMenuService(repo)
	handler := infrastructure.NewMenuHandler(service)

	// User menu query (needs to be registered outside the admin-restricted sub-group, but inside the jwt-protected router)
	router.GET("/menus/my", handler.GetMyMenus)

	menus := router.Group("/menus")
	{
		menus.GET("", middleware.RequirePermission(db, "/access-control/menus", "VIEW"), handler.GetAll)
		menus.GET("/:id", middleware.RequirePermission(db, "/access-control/menus", "VIEW"), handler.GetByID)
		menus.POST("", middleware.RequirePermission(db, "/access-control/menus", "CREATE"), handler.Create)
		menus.PUT("/:id", middleware.RequirePermission(db, "/access-control/menus", "EDIT"), handler.Update)
		menus.DELETE("/:id", middleware.RequirePermission(db, "/access-control/menus", "DELETE"), handler.Delete)
	}

	return repo
}
