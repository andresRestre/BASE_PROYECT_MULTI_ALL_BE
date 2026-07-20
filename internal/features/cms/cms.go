package cms

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"multicliente-backend/internal/features/cms/domain"
	"multicliente-backend/internal/platform/middleware"
)

func RegisterRoutes(publicGroup *gin.RouterGroup, protectedGroup *gin.RouterGroup, db *gorm.DB) {
	// -------------------------------------------------------------
	// 1. PUBLIC LANDING ENDPOINTS (For Next.js SSR)
	// -------------------------------------------------------------
	publicLanding := publicGroup.Group("/landing")
	{
		publicLanding.GET("/texts", func(c *gin.Context) {
			var texts []domain.LandingText
			if err := db.Order("id ASC").Find(&texts).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading texts"})
				return
			}
			c.JSON(http.StatusOK, texts)
		})

		publicLanding.GET("/news", func(c *gin.Context) {
			var newsList []domain.LandingNews
			if err := db.Where("is_published = ?", true).Order("create_at DESC").Find(&newsList).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading news"})
				return
			}
			c.JSON(http.StatusOK, newsList)
		})

		publicLanding.GET("/banners", func(c *gin.Context) {
			var banners []domain.LandingBanner
			if err := db.Where("is_active = ?", true).Order("sort_order ASC, id DESC").Find(&banners).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading banners"})
				return
			}
			c.JSON(http.StatusOK, banners)
		})
	}

	// -------------------------------------------------------------
	// 2. ADMIN PROTECTED CMS ENDPOINTS (For Flutter Admin Panel)
	// -------------------------------------------------------------
	cmsAdmin := protectedGroup.Group("/cms")
	{
		// --- TEXTS / ENCABEZADOS ---
		textsGroup := cmsAdmin.Group("/texts")
		{
			textsGroup.GET("", middleware.RequirePermission(db, "/website/texts", "VIEW"), func(c *gin.Context) {
				var texts []domain.LandingText
				if err := db.Order("id ASC").Find(&texts).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, texts)
			})

			textsGroup.PUT("/:id", middleware.RequirePermission(db, "/website/texts", "EDIT"), func(c *gin.Context) {
				id, _ := strconv.Atoi(c.Param("id"))
				var req domain.UpdateTextRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				var item domain.LandingText
				if err := db.First(&item, id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Text not found"})
					return
				}
				item.TextES = req.TextES
				item.TextEN = req.TextEN
				item.TextFR = req.TextFR
				db.Save(&item)
				c.JSON(http.StatusOK, item)
			})
		}

		// --- NOTICIAS ---
		newsGroup := cmsAdmin.Group("/news")
		{
			newsGroup.GET("", middleware.RequirePermission(db, "/website/news", "VIEW"), func(c *gin.Context) {
				var list []domain.LandingNews
				if err := db.Order("create_at DESC").Find(&list).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, list)
			})

			newsGroup.POST("", middleware.RequirePermission(db, "/website/news", "CREATE"), func(c *gin.Context) {
				var req domain.CreateNewsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				item := domain.LandingNews{
					TitleES:     req.TitleES,
					TitleEN:     req.TitleEN,
					TitleFR:     req.TitleFR,
					ContentES:   req.ContentES,
					ContentEN:   req.ContentEN,
					ContentFR:   req.ContentFR,
					ImageURL:    req.ImageURL,
					IsPublished: req.IsPublished,
				}
				if err := db.Create(&item).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, item)
			})

			newsGroup.PUT("/:id", middleware.RequirePermission(db, "/website/news", "EDIT"), func(c *gin.Context) {
				id, _ := strconv.Atoi(c.Param("id"))
				var req domain.UpdateNewsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				var item domain.LandingNews
				if err := db.First(&item, id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "News item not found"})
					return
				}
				if req.TitleES != nil { item.TitleES = *req.TitleES }
				if req.TitleEN != nil { item.TitleEN = *req.TitleEN }
				if req.TitleFR != nil { item.TitleFR = *req.TitleFR }
				if req.ContentES != nil { item.ContentES = *req.ContentES }
				if req.ContentEN != nil { item.ContentEN = *req.ContentEN }
				if req.ContentFR != nil { item.ContentFR = *req.ContentFR }
				if req.ImageURL != nil { item.ImageURL = *req.ImageURL }
				if req.IsPublished != nil { item.IsPublished = *req.IsPublished }
				db.Save(&item)
				c.JSON(http.StatusOK, item)
			})

			newsGroup.DELETE("/:id", middleware.RequirePermission(db, "/website/news", "DELETE"), func(c *gin.Context) {
				id, _ := strconv.Atoi(c.Param("id"))
				if err := db.Delete(&domain.LandingNews{}, id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
			})
		}

		// --- BANNERS DE IMÁGENES ---
		bannersGroup := cmsAdmin.Group("/banners")
		{
			bannersGroup.GET("", middleware.RequirePermission(db, "/website/banners", "VIEW"), func(c *gin.Context) {
				var list []domain.LandingBanner
				if err := db.Order("sort_order ASC, id DESC").Find(&list).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, list)
			})

			bannersGroup.POST("", middleware.RequirePermission(db, "/website/banners", "CREATE"), func(c *gin.Context) {
				var req domain.CreateBannerRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				item := domain.LandingBanner{
					Title:     req.Title,
					Subtitle:  req.Subtitle,
					ImageURL:  req.ImageURL,
					LinkURL:   req.LinkURL,
					SortOrder: req.SortOrder,
					IsActive:  req.IsActive,
				}
				if err := db.Create(&item).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, item)
			})

			bannersGroup.PUT("/:id", middleware.RequirePermission(db, "/website/banners", "EDIT"), func(c *gin.Context) {
				id, _ := strconv.Atoi(c.Param("id"))
				var req domain.UpdateBannerRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				var item domain.LandingBanner
				if err := db.First(&item, id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
					return
				}
				if req.Title != nil { item.Title = *req.Title }
				if req.Subtitle != nil { item.Subtitle = *req.Subtitle }
				if req.ImageURL != nil { item.ImageURL = *req.ImageURL }
				if req.LinkURL != nil { item.LinkURL = *req.LinkURL }
				if req.SortOrder != nil { item.SortOrder = *req.SortOrder }
				if req.IsActive != nil { item.IsActive = *req.IsActive }
				db.Save(&item)
				c.JSON(http.StatusOK, item)
			})

			bannersGroup.DELETE("/:id", middleware.RequirePermission(db, "/website/banners", "DELETE"), func(c *gin.Context) {
				id, _ := strconv.Atoi(c.Param("id"))
				if err := db.Delete(&domain.LandingBanner{}, id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Banner deleted successfully"})
			})
		}
	}
}
