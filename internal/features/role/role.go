package role

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/role/application"
	"multicliente-backend/internal/features/role/domain"
	"multicliente-backend/internal/features/role/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.RoleRepository {
	repo := infrastructure.NewRoleRepository(db)
	service := application.NewRoleService(repo)
	handler := infrastructure.NewRoleHandler(service)

	roles := router.Group("/roles")
	{
		roles.GET("", middleware.RequirePermission(db, "/roles", "VIEW"), handler.GetAll)
		roles.GET("/:id", middleware.RequirePermission(db, "/roles", "VIEW"), handler.GetByID)
		roles.GET("/options", middleware.RequirePermission(db, "/roles", "VIEW"), handler.GetOptions)
		
		roles.POST("", middleware.RequirePermission(db, "/roles", "CREATE"), handler.Create)
		roles.PUT("/:id", middleware.RequirePermission(db, "/roles", "EDIT"), handler.Update)
		roles.DELETE("/:id", middleware.RequirePermission(db, "/roles", "DELETE"), handler.Delete)
	}

	return repo
}
