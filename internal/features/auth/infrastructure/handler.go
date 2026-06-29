package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"multicliente-backend/internal/features/auth/domain"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile handles GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var userID uint
	if f, ok := userIDVal.(float64); ok {
		userID = uint(f)
	} else if u, ok := userIDVal.(uint); ok {
		userID = u
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format in session"})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
