package app

import (
	"gin-gorm-base/models"
	"gin-gorm-base/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C    *gin.Context
	uid  int
	user *models.User
}

type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"request_id"`
}

func (g *Gin) GetUser() *models.User {
	if nil == g.user {
		uid := g.GetUid()
		g.user, _ = models.GetUserById(uid)
	}
	return g.user
}

// GetUid 获取当前登录用户的 UID
func (g *Gin) GetUid() int {
	if g.uid <= 0 {
		g.uid = g.C.GetInt("uid")
	}
	return g.uid
}

// GetRoleId 获取当前登录用户的角色id
func (g *Gin) GetRoleId() int {
	if g.uid <= 0 {
		g.uid = g.C.GetInt("roleId")
	}
	return g.uid
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code:      errCode,
		Msg:       e.GetMsg(errCode),
		Data:      data,
		RequestId: g.GetRequestId(),
	})
	//logging.WithCtxDebug(g.C, "RequestURI"+g.C.Request.RequestURI, "data", data)
	return
}

func (g *Gin) GetRequestId() string {
	return g.C.Request.Header.Get("X-Request-ID")
}
