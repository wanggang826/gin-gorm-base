package v1

import (
	"gin-gorm-base/pkg/aliyun"
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OssStsInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	res := aliyun.GetStsToken()

	appG.Response(http.StatusOK, e.SUCCESS, res)
}
