package infrastructure

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"multicliente-backend/internal/features/auth/domain"
	"multicliente-backend/internal/platform/i18n"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	service domain.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	// Set httpOnly cookie with dynamic role session duration
	c.SetCookie("token", response.Token, response.SessionDurationSeconds, "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// GetProfile handles GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.TranslateError(c, errors.New("user not authenticated"))})
		return
	}

	var userID uint
	if f, ok := userIDVal.(float64); ok {
		userID = uint(f)
	} else if u, ok := userIDVal.(uint); ok {
		userID = u
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, errors.New("invalid user ID format in session"))})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword handles PUT /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n.TranslateError(c, errors.New("user not authenticated"))})
		return
	}

	var userID uint
	if f, ok := userIDVal.(float64); ok {
		userID = uint(f)
	} else if u, ok := userIDVal.(uint); ok {
		userID = u
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, errors.New("invalid user ID format"))})
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contraseña cambiada correctamente"})
}

// ForgotPassword handles POST /api/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	if err := h.service.ForgotPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Se ha enviado un correo con las instrucciones de recuperación"})
}

// ResetPassword handles POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.TranslateError(c, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña restablecida correctamente"})
}
