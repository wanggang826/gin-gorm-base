package models

import (
	"fmt"
	"gin-gorm-base/pkg/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	//"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/setting"
)

var db *gorm.DB

func Setup() {
	ConnStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)
	var gormLogger logger.Interface
	gormLogger = logging.NewGormLogger()
	if setting.DatabaseSetting.LogType == "output" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	engine, err := gorm.Open(mysql.New(
		mysql.Config{
			DSN:                       ConnStr,
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: gormLogger,
	})
	db = engine

	if err != nil {
		logging.Error("Fail to create gorm system logger: ", "msg", err.Error())
	} else {
		sqlDB, _ := db.DB()
		err := sqlDB.Ping()
		if err != nil {
			logging.Error("Connect to mysql error: ", "msg", err.Error())
		} else {
			//设置表明为单数形式
			//db.SingularTable(true)

			sqlDB.SetConnMaxLifetime(time.Hour)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetMaxIdleConns(10)

			logging.Info("Connect to sql OK : ", ConnStr)
		}
	}

}

func CheckFiled(paramsMap map[string]interface{}, dbFiled []string) bool {
	if len(dbFiled) == 0 {
		return false
	}
	for k := range paramsMap {
		if !util.IsInArray(k, dbFiled) {
			return false
		}
	}
	return true
}

type Writer struct {
}

// Printf 格式化外部组件日志
func (lw Writer) Printf(format string, v ...interface{}) {
	// 写入文件  这里就是你写入文件的地方
	logging.Info("db", "sql", v)
}
