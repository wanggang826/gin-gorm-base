package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set("X-Request-ID", requestID)
		}
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
