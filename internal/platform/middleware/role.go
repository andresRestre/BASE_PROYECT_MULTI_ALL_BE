package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a Gin middleware that checks if the authenticated user has one of the allowed roles.
// It relies on JWTAuth middleware having ran beforehand to inject role (role code) into context.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: role not found in context"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: invalid role format"})
			c.Abort()
			return
		}

		// Check if user's role is in the list of allowed roles
		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: insufficient permissions"})
		c.Abort()
	}
}
