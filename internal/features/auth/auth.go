package auth

import (
	"github.com/gin-gonic/gin"

	"multicliente-backend/internal/features/auth/application"
	"multicliente-backend/internal/features/auth/infrastructure"
	userDomain "multicliente-backend/internal/features/user/domain"
	"multicliente-backend/internal/platform/middleware"
)

// RegisterRoutes wires up the auth feature and registers its routes.
// It receives the user repository as a cross-feature dependency.
func RegisterRoutes(router *gin.RouterGroup, userRepo userDomain.UserRepository, jwtSecret string, jwtExpHours string) {
	service := application.NewAuthService(userRepo, jwtSecret, jwtExpHours)
	handler := infrastructure.NewAuthHandler(service)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.GET("/profile", middleware.JWTAuth(jwtSecret), handler.GetProfile)
		authGroup.PUT("/change-password", middleware.JWTAuth(jwtSecret), handler.ChangePassword)
	}
}
