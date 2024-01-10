package crontab

import (
	"gin-gorm-base/pkg/logging"
	"github.com/robfig/cron/v3"
)

func SetUp() {
	c := cron.New()
	//每天凌晨1点执行
	_, err := c.AddFunc("0 1 * * *", func() {
		logging.Debug("Cron SetUp cron.AddFunc ", "0 1")
	})
	if err != nil {
		logging.Error("Cron SetUp cron.AddFunc ", "err", err)
	}
	//每5秒执行一次
	//cron.AddFunc("@every s", func() {
	//	income_service.Test()
	//})
	c.Start()
}
