package v1

import (
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/service/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	result := user_service.GetUserInfo(appG.GetUid())
	appG.Response(http.StatusOK, e.SUCCESS, result)
}
