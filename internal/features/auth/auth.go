package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userDomain "multicliente-backend/internal/features/access_control/user/domain"
	"multicliente-backend/internal/features/auth/application"
	"multicliente-backend/internal/features/auth/infrastructure"
	"multicliente-backend/internal/platform/email"
	"multicliente-backend/internal/platform/middleware"
)

// RegisterRoutes wires up the auth feature and registers its routes.
func RegisterRoutes(
	router *gin.RouterGroup,
	userRepo userDomain.UserRepository,
	db *gorm.DB,
	emailService email.EmailService,
	jwtSecret string,
	jwtExpHours string,
) {
	service := application.NewAuthService(userRepo, db, emailService, jwtSecret, jwtExpHours)
	handler := infrastructure.NewAuthHandler(service)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/forgot-password", handler.ForgotPassword)
		authGroup.POST("/reset-password", handler.ResetPassword)

		authGroup.GET("/profile", middleware.JWTAuth(jwtSecret), handler.GetProfile)
		authGroup.PUT("/change-password", middleware.JWTAuth(jwtSecret), handler.ChangePassword)
	}
}
