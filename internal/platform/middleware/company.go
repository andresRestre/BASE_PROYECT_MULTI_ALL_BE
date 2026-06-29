package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireCompanyAccess checks if the X-Company-ID header is present and if the user is authorized.
// If valid, stores the company ID in the context as "active_company_id".
func RequireCompanyAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyIDStr := c.GetHeader("X-Company-ID")
		if companyIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Company-ID header is required"})
			c.Abort()
			return
		}

		companyIDVal, err := strconv.ParseUint(companyIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-Company-ID format"})
			c.Abort()
			return
		}
		companyID := uint(companyIDVal)

		// 1. Get role code from context
		roleCodeVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role information not found in token"})
			c.Abort()
			return
		}
		roleCode := roleCodeVal.(string)

		// SuperAdmin can access any company
		if roleCode == "superadmin" {
			c.Set("active_company_id", companyID)
			c.Next()
			return
		}

		// 2. Get user_id from context
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
			c.Abort()
			return
		}

		var userID uint
		if f, ok := userIDVal.(float64); ok {
			userID = uint(f)
		} else if u, ok := userIDVal.(uint); ok {
			userID = u
		} else if s, ok := userIDVal.(string); ok {
			parsed, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID in session"})
				c.Abort()
				return
			}
			userID = uint(parsed)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported user ID format"})
			c.Abort()
			return
		}

		// 3. Verify user is associated with this company in database
		var count int64
		err = db.Table("administrative.user_companies").
			Where("user_id = ? AND company_id = ?", userID, companyID).
			Count(&count).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify company permissions"})
			c.Abort()
			return
		}

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this company"})
			c.Abort()
			return
		}

		c.Set("active_company_id", companyID)
		c.Next()
	}
}
