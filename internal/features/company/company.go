package company

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/company/application"
	"multicliente-backend/internal/features/company/domain"
	"multicliente-backend/internal/features/company/infrastructure"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, requireRole gin.HandlerFunc) domain.CompanyRepository {
	repo := infrastructure.NewCompanyRepository(db)
	service := application.NewCompanyService(repo)
	handler := infrastructure.NewCompanyHandler(service)

	companies := router.Group("/companies")
	{
		companies.GET("", middleware.RequirePermission(db, "/companies", "VIEW"), handler.GetAll)
		companies.GET("/:id", middleware.RequirePermission(db, "/companies", "VIEW"), handler.GetByID)
		
		companies.POST("", middleware.RequirePermission(db, "/companies", "CREATE"), handler.Create)
		companies.PUT("/:id", middleware.RequirePermission(db, "/companies", "EDIT"), handler.Update)
		companies.DELETE("/:id", middleware.RequirePermission(db, "/companies", "DELETE"), handler.Delete)
	}

	return repo
}
