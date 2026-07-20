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

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	companyIDVal, exists := c.Get("companyID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company context missing"})
		return
	}
	companyID := companyIDVal.(uint)

	list, err := h.service.GetNotifications(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userIDVal.(uint)

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
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	companyIDVal, exists := c.Get("companyID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company context missing"})
		return
	}
	companyID := companyIDVal.(uint)

	err := h.service.MarkAllRead(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
