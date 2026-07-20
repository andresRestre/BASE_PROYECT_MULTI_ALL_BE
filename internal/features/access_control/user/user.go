package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/access_control/user/application"
	"multicliente-backend/internal/features/access_control/user/domain"
	"multicliente-backend/internal/features/access_control/user/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

// RegisterRoutes wires up the user feature: repository → service → handler → routes.
// Returns the UserRepository so other features (e.g., auth) can reuse it.
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) domain.UserRepository {
	repo := infrastructure.NewUserRepository(db)
	service := application.NewUserService(repo)
	handler := infrastructure.NewUserHandler(service)

	users := router.Group("/users")
	{
		users.POST("", middleware.RequirePermission(db, "/access-control/users", "CREATE"), handler.Create)
		users.GET("", middleware.RequirePermission(db, "/access-control/users", "VIEW"), handler.GetAll)
		users.GET("/:id", middleware.RequirePermission(db, "/access-control/users", "VIEW"), handler.GetByID)
		users.PUT("/:id", middleware.RequirePermission(db, "/access-control/users", "EDIT"), handler.Update)
		users.DELETE("/:id", middleware.RequirePermission(db, "/access-control/users", "DELETE"), handler.Delete)
	}

	return repo
}
