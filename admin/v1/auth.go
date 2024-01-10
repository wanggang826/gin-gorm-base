package v1

import (
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/service/admin_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 登录
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	token, code := admin_service.AdminLogin(c)
	appG.Response(http.StatusOK, code, token)
}
