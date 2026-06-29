package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"multicliente-backend/internal/features/menu/domain"
)

type MenuHandler struct {
	service domain.MenuService
}

func NewMenuHandler(service domain.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req domain.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := getUserIDFromContext(c)

	menu, err := h.service.CreateMenu(&req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}
	id := uint(idVal)

	menu, err := h.service.GetMenuByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}
	id := uint(idVal)

	var req domain.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBy := getUserIDFromContext(c)

	menu, err := h.service.UpdateMenu(id, &req, updatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu ID"})
		return
	}
	id := uint(idVal)

	if err := h.service.DeleteMenu(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "menu deleted successfully"})
}

func (h *MenuHandler) GetMyMenus(c *gin.Context) {
	roleIDVal, exists := c.Get("role_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role information not found in token"})
		return
	}

	var roleID uint
	if f, ok := roleIDVal.(float64); ok {
		roleID = uint(f)
	} else if u, ok := roleIDVal.(uint); ok {
		roleID = u
	} else if s, ok := roleIDVal.(string); ok {
		parsed, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID in token"})
			return
		}
		roleID = uint(parsed)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported role ID format"})
		return
	}

	allowedMenus, err := h.service.GetAllowedMenus(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allowedMenus)
}

func getUserIDFromContext(c *gin.Context) *uint {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return nil
	}

	if f, ok := userIDVal.(float64); ok {
		u := uint(f)
		return &u
	}

	if s, ok := userIDVal.(string); ok {
		idVal, err := strconv.ParseUint(s, 10, 32)
		if err == nil {
			u := uint(idVal)
			return &u
		}
	}

	return nil
}
