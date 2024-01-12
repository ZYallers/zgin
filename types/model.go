package types

import (
	"time"

	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/golib/utils/mysql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Model struct {
	mysql.Model
}

func (m *Model) Config(db string, options ...interface{}) func() *gorm.Config {
	logLevel := gormLogger.Warn
	slowThreshold := 500 * time.Millisecond
	ol := len(options)
	if ol == 0 && gin.IsDebugging() {
		logLevel = gormLogger.Info
	}
	if ol > 0 {
		logLevel = options[0].(gormLogger.LogLevel)
	}
	if ol > 1 {
		slowThreshold = options[1].(time.Duration)
	}
	return func() *gorm.Config {
		return &gorm.Config{Logger: logger.NewGormLogger(db, slowThreshold, logLevel, &Sender{})}
	}
}
