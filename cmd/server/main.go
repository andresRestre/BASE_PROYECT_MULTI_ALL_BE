package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"multicliente-backend/internal/features/access_control"
	articleDomain "multicliente-backend/internal/features/inventory/article/domain"
	"multicliente-backend/internal/features/auth"
	authDomain "multicliente-backend/internal/features/auth/domain"
	categoryDomain "multicliente-backend/internal/features/inventory/category/domain"
	companyDomain "multicliente-backend/internal/features/access_control/company/domain"
	"multicliente-backend/internal/features/cms"
	cmsDomain "multicliente-backend/internal/features/cms/domain"
	menuDomain "multicliente-backend/internal/features/access_control/menu/domain"
	roleDomain "multicliente-backend/internal/features/access_control/role/domain"
	"multicliente-backend/internal/features/inventory"
	notificationDomain "multicliente-backend/internal/features/notification/domain"
	"multicliente-backend/internal/features/notification"
	"multicliente-backend/internal/features/upload"
	userDomain "multicliente-backend/internal/features/access_control/user/domain"
	"multicliente-backend/internal/platform/config"
	"multicliente-backend/internal/platform/database"
	"multicliente-backend/internal/platform/database/migrations"
	"multicliente-backend/internal/platform/database/seeds"
	"multicliente-backend/internal/platform/email"
	"multicliente-backend/internal/platform/middleware"
	"multicliente-backend/internal/platform/server"
)

func main() {
	// Load configuration from .env
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// Run auto-migrations in order of dependency
	// Note: Drop legacy UUID columns if they exist in administrative.users, as PostgreSQL cannot alter UUID to bigint directly
	_ = db.Exec("ALTER TABLE administrative.users DROP COLUMN IF EXISTS create_by")
	_ = db.Exec("ALTER TABLE administrative.users DROP COLUMN IF EXISTS update_by")

	err = migrations.Migrate(db,
		&companyDomain.Company{},
		&categoryDomain.Category{},
		&articleDomain.Article{},
		&roleDomain.Role{},
		&roleDomain.Option{},
		&roleDomain.Permission{},
		&roleDomain.RoleNotificationRule{},
		&menuDomain.Menu{},
		&userDomain.User{},
		&authDomain.PasswordResetToken{},
		&notificationDomain.Notification{},
		&cmsDomain.LandingText{},
		&cmsDomain.LandingNews{},
		&cmsDomain.LandingBanner{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Seed database default values
	seeds.Seed(db)

	// Initialize services
	emailService := email.NewEmailService(cfg)

	// Setup Gin router
	router := server.NewRouter()

	// API route group
	api := router.Group("/api")

	// Protected routes (JWT required)
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))

	// Middleware to restrict endpoints to SuperAdmin role only
	superAdminRequired := middleware.RequireRole("superadmin")
	requireCompanyAccess := middleware.RequireCompanyAccess(db)

	// Register features
	accessControlRouter := protected.Group("/access-control")
	userRepo := access_control.RegisterRoutes(accessControlRouter, db, superAdminRequired)
	auth.RegisterRoutes(api, userRepo, db, emailService, cfg.JWTSecret, cfg.JWTExpirationHours)

	// Register notifications
	notificationService := notification.RegisterRoutes(protected, db)
	
	inventoryRouter := protected.Group("/inventory")
	inventory.RegisterRoutes(inventoryRouter, db, requireCompanyAccess, notificationService)

	upload.RegisterRoutes(protected)

	// Register CMS / Landing Page Routes (Public & Protected)
	cms.RegisterRoutes(api, protected, db)

	// Health check (public)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
