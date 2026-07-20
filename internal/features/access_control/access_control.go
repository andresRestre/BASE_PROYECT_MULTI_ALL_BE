package access_control

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"multicliente-backend/internal/features/access_control/company"
	"multicliente-backend/internal/features/access_control/menu"
	"multicliente-backend/internal/features/access_control/role"
	"multicliente-backend/internal/features/access_control/user"
	userDomain "multicliente-backend/internal/features/access_control/user/domain"
)

// RegisterRoutes registers all access control feature routes (user, company, role, menu) under a parent route group.
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, superAdminRequired gin.HandlerFunc) userDomain.UserRepository {
	userRepo := user.RegisterRoutes(router, db)
	company.RegisterRoutes(router, db, superAdminRequired)
	role.RegisterRoutes(router, db, superAdminRequired)
	menu.RegisterRoutes(router, db, superAdminRequired)
	return userRepo
}
