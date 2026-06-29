package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"multicliente-backend/internal/platform/localization"
)

type translationWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *translationWriter) WriteHeader(status int) {
	w.status = status
}

func (w *translationWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *translationWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *translationWriter) Status() int {
	if w.status == 0 {
		return w.ResponseWriter.Status()
	}
	return w.status
}

// Translation intercepts JSON error responses and translates the error message based on the Accept-Language header.
func Translation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tw := &translationWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
			status:         0,
		}
		c.Writer = tw

		c.Next()

		lang := c.GetHeader("Accept-Language")
		contentType := tw.Header().Get("Content-Type")

		// Retrieve the status code to write
		status := tw.status
		if status == 0 {
			status = http.StatusOK
		}

		if strings.Contains(contentType, "application/json") && tw.body.Len() > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal(tw.body.Bytes(), &data); err == nil {
				if errMsg, exists := data["error"].(string); exists {
					// Translate error message
					data["error"] = localization.Translate(lang, errMsg)
					
					// Marshal the updated response body
					newBody, err := json.Marshal(data)
					if err == nil {
						tw.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(newBody)))
						tw.ResponseWriter.WriteHeader(status)
						tw.ResponseWriter.Write(newBody)
						return
					}
				}
			}
		}

		// Fallback: Write the original status and body if not JSON or if it has no error message
		if tw.status != 0 {
			tw.ResponseWriter.WriteHeader(tw.status)
		}
		tw.ResponseWriter.Write(tw.body.Bytes())
	}
}
