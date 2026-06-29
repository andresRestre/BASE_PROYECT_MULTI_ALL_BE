package role

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/role/application"
	"multicliente-backend/internal/features/role/domain"
	"multicliente-backend/internal/features/role/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.RoleRepository {
	repo := infrastructure.NewRoleRepository(db)
	service := application.NewRoleService(repo)
	handler := infrastructure.NewRoleHandler(service)

	roles := router.Group("/roles")
	{
		roles.GET("", handler.GetAll)
		roles.GET("/:id", handler.GetByID)
		roles.GET("/options", handler.GetOptions)
		
		// SuperAdmin only endpoints
		roles.POST("", requireRole, handler.Create)
		roles.PUT("/:id", requireRole, handler.Update)
		roles.DELETE("/:id", requireRole, handler.Delete)
	}

	return repo
}
