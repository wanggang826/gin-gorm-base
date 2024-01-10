package routers

import (
	v1 "gin-gorm-base/admin/v1"
	"gin-gorm-base/middleware/jwt"
)

func InitAdminRouter() {
	adminV1 := r.Group("/admin/v1")
	adminV1.POST("/auth/login", v1.Login)
	adminV1.GET("/comm/ossStsInfo", v1.OssStsInfo)
	adminV1.Use(jwt.JWT())
	{
		// admin 账号
		adminV1.GET("/admin/getPage", v1.GetAdminPage)
		adminV1.GET("/admin/search", v1.SearchAdmin)
		adminV1.POST("/admin/add", v1.AddAdmin)
		adminV1.POST("/admin/edit", v1.EditAdmin)
		adminV1.POST("/admin/del", v1.DelAdmin)
		adminV1.POST("/admin/changeStatus", v1.ChangeAdminStatus)
		adminV1.POST("/admin/changePwd", v1.EditAdminPwd)
		adminV1.POST("/admin/resetPwd", v1.ResetAdminPwd)
		adminV1.GET("/admin/info", v1.GetAdminInfo)

		// notice 公告
		adminV1.GET("/notice/getPage", v1.GetNoticePage)
		adminV1.POST("/notice/add", v1.AddNotice)
		adminV1.POST("/notice/edit", v1.EditNotice)
		adminV1.POST("/notice/del", v1.DelNotice)

	}
}
