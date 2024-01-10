package v1

import (
	"gin-gorm-base/pkg/douyin"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/service/user_service"
	"github.com/gin-gonic/gin"
	"net/http"

	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/pkg/util"
)

type LoginParams struct {
	Code          string `json:"code"`
	AnonymousCode string `json:"anonymous_code"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var params LoginParams
	err := c.BindJSON(&params)
	if err != nil {
		logging.Error("GetRequest: Login Gin BindJSON:", "msg", err.Error())
		appG.Response(http.StatusOK, e.ERROR_PARAM_GIN_BINDJSON_ERROR, nil)
		return
	}
	userId, err := user_service.SaveUser(params.Code, params.AnonymousCode)
	if err != nil {
		logging.Error("UpdateUser user info  error:", err)
		appG.Response(http.StatusOK, e.ERROR_LOGIN_FAIL, nil)
		return
	}
	token, err := util.GenerateUserToken(userId, 0)
	if err != nil {
		logging.Error("GenerateUserToken:", "userId", userId, "token", token)
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	openId, anonymousOpenid := user_service.GetUserOpenId(userId)
	resp := make(map[string]interface{})
	resp["token"] = token
	resp["anonymous_openid"] = anonymousOpenid
	resp["openid"] = openId
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

func RefreshToken(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := appG.GetUid()

	token, err := util.GenerateUserToken(userId, 0)
	if err != nil {
		logging.Error("RefreshToken err:", "msg", err.Error(), "userId", userId)
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, token)
}

func TestGetAccessToken(c *gin.Context) {
	appG := app.Gin{C: c}
	dyService := douyin.NewDyService()
	accessToken := dyService.GetAccessToken()
	result := make(map[string]interface{})
	result["accessToken"] = accessToken
	appG.Response(http.StatusOK, e.SUCCESS, result)
}
