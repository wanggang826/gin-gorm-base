package main

import (
	"fmt"
	"gin-gorm-base/crontab"
	"gin-gorm-base/pkg/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"gin-gorm-base/models"
	"gin-gorm-base/pkg/gredis"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/setting"
	"gin-gorm-base/routers"
)

func init() {
	gin.SetMode(setting.ServerSetting.RunMode)
	setting.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
	models.Setup()
	crontab.SetUp()
}

func main() {
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	server := &http.Server{
		Addr:         endPoint,
		Handler:      routersInit,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
