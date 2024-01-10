package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-gorm-base/pkg/e"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code = e.SUCCESS
		var data interface{}

		uid, roleId, isAdmin, err := util.GetUidFromHeader(c)
		logging.Debug("JWT GetUidFromHeader res", "uid", uid, "isAdmin", isAdmin, "err", err)
		if err != nil {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			logging.Error("jwt check token error: ", "msg", err.Error(), "uid", uid, "isAdmin", isAdmin)
		} else {
			if !isAdmin && c.Request.URL.Path[:7] == "/admin/" {
				logging.Error("jwt ILLEGAL request : ", uid)
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else {
				c.Set("uid", uid)
				c.Set("roleId", roleId)
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
