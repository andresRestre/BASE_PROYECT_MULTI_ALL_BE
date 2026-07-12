package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"multicliente-backend/internal/features/category/domain"
	"multicliente-backend/internal/platform/i18n"
)

type CategoryHandler struct {
	service domain.CategoryService
}

func NewCategoryHandler(service domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}
func (h *CategoryHandler) Create(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	companyID := c.MustGet("active_company_id").(uint)
	createdBy := getUserIDFromContext(c)

	cat, err := h.service.CreateCategory(&req, companyID, createdBy)
	if err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	companyID := c.MustGet("active_company_id").(uint)

	categories, err := h.service.GetCategoriesByCompany(companyID)
	if err != nil {
		i18n.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid category ID")
		return
	}
	id := uint(idVal)

	cat, err := h.service.GetCategoryByID(id)
	if err != nil {
		i18n.Error(c, http.StatusNotFound, err)
		return
	}

	// Verify the category belongs to the active company context
	companyID := c.MustGet("active_company_id").(uint)
	if cat.CompanyID != companyID {
		i18n.ErrorString(c, http.StatusForbidden, "access_denied")
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid category ID")
		return
	}
	id := uint(idVal)

	// Verify active company context matches category
	cat, err := h.service.GetCategoryByID(id)
	if err != nil {
		i18n.Error(c, http.StatusNotFound, err)
		return
	}
	companyID := c.MustGet("active_company_id").(uint)
	if cat.CompanyID != companyID {
		i18n.ErrorString(c, http.StatusForbidden, "access_denied")
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	updatedBy := getUserIDFromContext(c)

	updatedCat, err := h.service.UpdateCategory(id, &req, updatedBy)
	if err != nil {
		i18n.Error(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, updatedCat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		i18n.ErrorString(c, http.StatusBadRequest, "invalid category ID")
		return
	}
	id := uint(idVal)

	// Verify active company context matches category
	cat, err := h.service.GetCategoryByID(id)
	if err != nil {
		i18n.Error(c, http.StatusNotFound, err)
		return
	}
	companyID := c.MustGet("active_company_id").(uint)
	if cat.CompanyID != companyID {
		i18n.ErrorString(c, http.StatusForbidden, "access_denied")
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		i18n.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
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

	if u, ok := userIDVal.(uint); ok {
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
