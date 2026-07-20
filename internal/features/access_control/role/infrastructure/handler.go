package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"multicliente-backend/internal/features/access_control/role/domain"
	"multicliente-backend/internal/platform/i18n"
)

type RoleHandler struct {
	service domain.RoleService
}

func NewRoleHandler(service domain.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req domain.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	createdBy := getUserIDFromContext(c)

	role, err := h.service.CreateRole(&req, createdBy)
	if err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, role)
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		i18n.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid role ID")
		return
	}
	id := uint(idVal)

	role, err := h.service.GetRoleByID(id)
	if err != nil {
		i18n.Error(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid role ID")
		return
	}
	id := uint(idVal)

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	updatedBy := getUserIDFromContext(c)

	role, err := h.service.UpdateRole(id, &req, updatedBy)
	if err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid role ID")
		return
	}
	id := uint(idVal)
	deletedBy := getUserIDFromContext(c)

	if err := h.service.DeleteRole(id, deletedBy); err != nil {
		i18n.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

func (h *RoleHandler) GetOptions(c *gin.Context) {
	options, err := h.service.GetAllOptions()
	if err != nil {
		i18n.Error(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, options)
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
