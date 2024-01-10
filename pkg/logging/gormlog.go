package logging

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	LogLevel                            logger.LogLevel
	SlowThreshold                       time.Duration
	IgnoreRecordNotFoundError           bool
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormLogger() *GormLogger {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &GormLogger{
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
		LogLevel:     logger.Info,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		Logger.Info(msg, formatParams(data...)...)
	}
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		Logger.Warn(msg, formatParams(data...)...)
	}
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		Logger.Error(msg, formatParams(data...)...)
	}
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			Error("ERROR", "line", utils.FileWithLineNum(), "error", err.Error(),
				"time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			Error("ERROR", "line", utils.FileWithLineNum(), "error", err.Error(),
				"time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			Warn("WARN", "line", utils.FileWithLineNum(), "slowLog", slowLog,
				"time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			Warn("WARN", "line", utils.FileWithLineNum(), "slowLog", slowLog,
				"time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			Info("Info", "line", utils.FileWithLineNum(), "time", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			Info("Info", "line", utils.FileWithLineNum(), "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}
