package routers

import (
	v1 "gin-gorm-base/api/v1"
	"gin-gorm-base/middleware/jwt"
)

func InitApiRouter() {
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/auth/testGetAccessToken", v1.TestGetAccessToken)
	apiV1.POST("/auth/login", v1.Login)
	apiV1.Use(jwt.JWT())
	{
		//refresh token
		apiV1.GET("/auth/refreshToken", v1.RefreshToken)

		//user
		apiV1.GET("/user/get", v1.GetUserInfo)

	}
}
