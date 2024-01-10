package v1

import (
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/service/admin_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAdminPage 列表查询
func GetAdminPage(c *gin.Context) {
	appG := app.Gin{C: c}
	result := admin_service.GetAdminPage(c)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

// GetAdminInfo 账号信息
func GetAdminInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	result := admin_service.GetAdminInfo(c)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

// SearchAdmin 搜索
func SearchAdmin(c *gin.Context) {
	appG := app.Gin{C: c}
	result := admin_service.SearchAdmin(c)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

// AddAdmin 新增
func AddAdmin(c *gin.Context) {
	appG := app.Gin{C: c}
	result, code := admin_service.AddAdmin(c)
	appG.Response(http.StatusOK, code, result)
}

// EditAdmin 编辑
func EditAdmin(c *gin.Context) {
	appG := app.Gin{C: c}
	result, code := admin_service.EditAdmin(c)
	appG.Response(http.StatusOK, code, result)
}

// DelAdmin 删除
func DelAdmin(c *gin.Context) {
	appG := app.Gin{C: c}
	code := admin_service.DelAdmin(c)
	appG.Response(http.StatusOK, code, nil)
}

// ChangeAdminStatus 启用|禁用
func ChangeAdminStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	code := admin_service.ChangeAdminStatus(c)
	appG.Response(http.StatusOK, code, nil)
}

// EditAdminPwd 修改密码
func EditAdminPwd(c *gin.Context) {
	appG := app.Gin{C: c}
	code := admin_service.EditAdminPwd(c)
	appG.Response(http.StatusOK, code, nil)
}

// ResetAdminPwd 重置密码
func ResetAdminPwd(c *gin.Context) {
	appG := app.Gin{C: c}
	code := admin_service.ResetAdminPwd(c)
	appG.Response(http.StatusOK, code, nil)
}
