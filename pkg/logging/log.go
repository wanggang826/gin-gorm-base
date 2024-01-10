package logging

import (
	"fmt"
	"gin-gorm-base/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Setup() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	logLevel := zapcore.ErrorLevel
	if setting.ServerSetting.RunMode == "debug" {
		logLevel = zapcore.DebugLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Logger.Info("zap.Logger inited")
}

// Debug output logs at debug level
func Debug(msg string, params ...interface{}) {
	Logger.Debug(msg, formatParams(params...)...)
}

// Info output logs at info level
func Info(msg string, params ...interface{}) {
	Logger.Info(msg, formatParams(params...)...)
}

// Warn output logs at warn level
func Warn(msg string, params ...interface{}) {
	Logger.Warn(msg, formatParams(params...)...)
}

// Error output logs at error level
func Error(msg string, params ...interface{}) {
	fields := formatParams(params...)
	Logger.Error(msg, fields...)
}

// Fatal output logs at fatal level
func Fatal(msg string, params ...interface{}) {
	Logger.Fatal(msg, formatParams(params...)...)
}

// WithCtxDebug output logs at debug level
func WithCtxDebug(c *gin.Context, msg string, params ...interface{}) {
	params = append(params, "Request-ID", c.Request.Header.Get("X-Request-ID"))
	Logger.Debug(msg, formatParams(params...)...)
}

// WithCtxInfo output logs at info level
func WithCtxInfo(c *gin.Context, msg string, params ...interface{}) {
	params = append(params, "Request-ID", c.Request.Header.Get("X-Request-ID"))
	Logger.Info(msg, formatParams(params...)...)
}

// WithCtxWarn output logs at warn level
func WithCtxWarn(c *gin.Context, msg string, params ...interface{}) {
	params = append(params, "Request-ID", c.Request.Header.Get("X-Request-ID"))
	Logger.Warn(msg, formatParams(params...)...)
}

// WithCtxError output logs at error level
func WithCtxError(c *gin.Context, msg string, params ...interface{}) {
	params = append(params, "Request-ID", c.Request.Header.Get("X-Request-ID"))
	fields := formatParams(params...)
	Logger.Error(msg, fields...)
}

// WithCtxFatal output logs at fatal level
func WithCtxFatal(c *gin.Context, msg string, params ...interface{}) {
	params = append(params, "Request-ID", c.Request.Header.Get("X-Request-ID"))
	Logger.Fatal(msg, formatParams(params...)...)
}

func formatParams(params ...interface{}) []zapcore.Field {
	total := len(params)

	//把参数填成偶数个
	if total%2 == 1 {
		params = append(params, "")
		total += 1
	}

	var fields []zapcore.Field
	for i := 0; i < total; i += 2 {
		field := zap.Any(fmt.Sprintf("%+v", params[i]), params[i+1])
		fields = append(fields, field)
	}
	return fields
}

func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogWriter() zapcore.WriteSyncer {
	filePath := getLogFilePath()

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./" + filePath + "latest.log", // 日志输出文件
		MaxSize:    100,                            // 日志最大保存100M
		MaxBackups: 3,                              // 旧日志保留7个日志备份
		MaxAge:     7,                              // 日志存活时长 最多保留7天日志
		Compress:   true,                           // 自导打 gzip包 默认false
	}
	return zapcore.AddSync(lumberJackLogger)
}
