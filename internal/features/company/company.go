package company

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/company/application"
	"multicliente-backend/internal/features/company/domain"
	"multicliente-backend/internal/features/company/infrastructure"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.CompanyRepository {
	repo := infrastructure.NewCompanyRepository(db)
	service := application.NewCompanyService(repo)
	handler := infrastructure.NewCompanyHandler(service)

	companies := router.Group("/companies")
	{
		companies.GET("", handler.GetAll)
		companies.GET("/:id", handler.GetByID)
		
		// SuperAdmin only endpoints
		companies.POST("", requireRole, handler.Create)
		companies.PUT("/:id", requireRole, handler.Update)
		companies.DELETE("/:id", requireRole, handler.Delete)
	}

	return repo
}
