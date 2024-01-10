package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	PrefixUrl       string
	RuntimeRootPath string
	LogSavePath     string
	JwtUserSecret   string
	JwtAdminSecret  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
	LogType     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	Db          int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Oss struct {
	OssDomain          string
	AccessKeyId        string
	AccessKeySecret    string
	OssEndpoint        string
	OssDefaultBucket   string
	OssRegion          string
	OssRoleARN         string
	OssRoleSessionName string
}

var OssSetting = &Oss{}

type DouYin struct {
	Appid           string
	Secret          string
	Salt            string
	Token           string
	NotifyUrl       string
	RefundNotifyUrl string
}

var DouYinSetting = &DouYin{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")

	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("oss", OssSetting)
	mapTo("douYin", DouYinSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
