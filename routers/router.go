package routers

import (
	"gin-gorm-base/middleware/requestid"
	"gin-gorm-base/pkg/logging"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

var r *gin.Engine

func InitRouter() *gin.Engine {

	r = gin.New()
	r.Use(ginzap.Ginzap(logging.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logging.Logger, true))

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to gin-gorm-base1")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	r.Use(requestid.GenerateRequestID())

	InitApiRouter()
	InitAdminRouter()
	return r
}
