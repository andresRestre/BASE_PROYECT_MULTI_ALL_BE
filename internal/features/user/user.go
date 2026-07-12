package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/user/application"
	"multicliente-backend/internal/features/user/domain"
	"multicliente-backend/internal/features/user/infrastructure"
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
		users.POST("", middleware.RequirePermission(db, "/users", "CREATE"), handler.Create)
		users.GET("", middleware.RequirePermission(db, "/users", "VIEW"), handler.GetAll)
		users.GET("/:id", middleware.RequirePermission(db, "/users", "VIEW"), handler.GetByID)
		users.PUT("/:id", middleware.RequirePermission(db, "/users", "EDIT"), handler.Update)
		users.DELETE("/:id", middleware.RequirePermission(db, "/users", "DELETE"), handler.Delete)
	}

	return repo
}
