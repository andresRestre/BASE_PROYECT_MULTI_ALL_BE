package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the upload endpoint and ensures the uploads directory exists.
func RegisterRoutes(router *gin.RouterGroup) {
	// Create uploads directory if it does not exist
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		fmt.Printf("❌ Failed to create uploads directory: %v\n", err)
	}

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Save the file with a unique name based on timestamp to avoid conflicts
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("./uploads", newFilename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Return the access path
		c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("/uploads/%s", newFilename)})
	})
}
