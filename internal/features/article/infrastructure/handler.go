package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"multicliente-backend/internal/features/article/domain"
)

type ArticleHandler struct {
	service domain.ArticleService
}

func NewArticleHandler(service domain.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req domain.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID := c.MustGet("active_company_id").(uint)
	createdBy := getUserIDFromContext(c)

	art, err := h.service.CreateArticle(&req, companyID, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, art)
}

func (h *ArticleHandler) GetAll(c *gin.Context) {
	companyID := c.MustGet("active_company_id").(uint)

	articles, err := h.service.GetArticlesByCompany(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetByID(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}
	id := uint(idVal)

	art, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	companyID := c.MustGet("active_company_id").(uint)
	if art.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this article context"})
		return
	}

	c.JSON(http.StatusOK, art)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}
	id := uint(idVal)

	art, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("active_company_id").(uint)
	if art.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this article context"})
		return
	}

	var req domain.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBy := getUserIDFromContext(c)

	updatedArt, err := h.service.UpdateArticle(id, &req, updatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedArt)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	idVal, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article ID"})
		return
	}
	id := uint(idVal)

	art, err := h.service.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("active_company_id").(uint)
	if art.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this article context"})
		return
	}

	if err := h.service.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
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
