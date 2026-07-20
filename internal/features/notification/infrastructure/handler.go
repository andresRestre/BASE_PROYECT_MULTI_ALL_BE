package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"multicliente-backend/internal/features/notification/domain"
)

type NotificationHandler struct {
	service domain.NotificationService
}

func NewNotificationHandler(service domain.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists || val == nil {
		val, exists = c.Get("userID")
	}
	if !exists || val == nil {
		return 0, false
	}
	if f, ok := val.(float64); ok {
		return uint(f), true
	}
	if u, ok := val.(uint); ok {
		return u, true
	}
	if p, ok := val.(*uint); ok && p != nil {
		return *p, true
	}
	return 0, false
}

func getCompanyID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("company_id")
	if !exists || val == nil {
		val, exists = c.Get("companyID")
	}
	if exists && val != nil {
		if f, ok := val.(float64); ok {
			return uint(f), true
		}
		if u, ok := val.(uint); ok {
			return u, true
		}
		if p, ok := val.(*uint); ok && p != nil {
			return *p, true
		}
	}
	// Fallback to X-Company-ID header
	header := c.GetHeader("X-Company-ID")
	if header != "" {
		if parsed, err := strconv.ParseUint(header, 10, 32); err == nil && parsed > 0 {
			return uint(parsed), true
		}
	}
	return 0, false
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	companyID, ok := getCompanyID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company context missing"})
		return
	}

	list, err := h.service.GetNotifications(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}
	id := uint(idVal)

	notif, err := h.service.MarkAsRead(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notif)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	companyID, ok := getCompanyID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company context missing"})
		return
	}

	err := h.service.MarkAllRead(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
