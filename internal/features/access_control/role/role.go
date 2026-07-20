package role

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/access_control/role/application"
	"multicliente-backend/internal/features/access_control/role/domain"
	"multicliente-backend/internal/features/access_control/role/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.RoleRepository {
	repo := infrastructure.NewRoleRepository(db)
	service := application.NewRoleService(repo)
	handler := infrastructure.NewRoleHandler(service)

	roles := router.Group("/roles")
	{
		roles.GET("", middleware.RequirePermission(db, "/access-control/roles", "VIEW"), handler.GetAll)
		roles.GET("/:id", middleware.RequirePermission(db, "/access-control/roles", "VIEW"), handler.GetByID)
		roles.GET("/options", middleware.RequirePermission(db, "/access-control/roles", "VIEW"), handler.GetOptions)
		
		roles.POST("", middleware.RequirePermission(db, "/access-control/roles", "CREATE"), handler.Create)
		roles.PUT("/:id", middleware.RequirePermission(db, "/access-control/roles", "EDIT"), handler.Update)
		roles.DELETE("/:id", middleware.RequirePermission(db, "/access-control/roles", "DELETE"), handler.Delete)
	}

	return repo
}
