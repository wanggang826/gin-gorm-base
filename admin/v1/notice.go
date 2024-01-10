package v1

import (
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/service/notice_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetNoticePage 列表查询
func GetNoticePage(c *gin.Context) {
	appG := app.Gin{C: c}
	result := notice_service.GetNoticePage(c)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

// AddNotice 新增
func AddNotice(c *gin.Context) {
	appG := app.Gin{C: c}
	result, code := notice_service.AddNotice(c)
	appG.Response(http.StatusOK, code, result)
}

// EditNotice 编辑
func EditNotice(c *gin.Context) {
	appG := app.Gin{C: c}
	result, code := notice_service.EditNotice(c)
	appG.Response(http.StatusOK, code, result)
}

// DelNotice 删除
func DelNotice(c *gin.Context) {
	appG := app.Gin{C: c}
	code := notice_service.DelNotice(c)
	appG.Response(http.StatusOK, code, nil)
}
