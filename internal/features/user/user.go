package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/user/application"
	"multicliente-backend/internal/features/user/domain"
	"multicliente-backend/internal/features/user/infrastructure"
)

// RegisterRoutes wires up the user feature: repository → service → handler → routes.
// Returns the UserRepository so other features (e.g., auth) can reuse it.
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) domain.UserRepository {
	repo := infrastructure.NewUserRepository(db)
	service := application.NewUserService(repo)
	handler := infrastructure.NewUserHandler(service)

	users := router.Group("/users")
	{
		users.POST("", handler.Create)
		users.GET("", handler.GetAll)
		users.GET("/:id", handler.GetByID)
		users.PUT("/:id", handler.Update)
		users.DELETE("/:id", handler.Delete)
	}

	return repo
}
