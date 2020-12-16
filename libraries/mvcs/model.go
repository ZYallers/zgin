package mvcs

import (
	"errors"
	"fmt"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

const (
	defaultCharset          = "utf8mb4"
	defaultLoc              = "Local"
	defaultParseTime        = "true"
	defaultMaxAllowedPacket = "0"
	defaultTimeout          = "15s"
	retryConnDbMaxTimes     = 3
	maxOpenConns            = 1000
)

type Model struct {
}

type DbCollector struct {
	done    uint32
	pointer *gorm.DB
}

func (m *Model) NewClient(dbc *DbCollector, dialect *app.MysqlDialect) *gorm.DB {
	defer tool.SafeDefer()
	var err error
	for i := 1; i <= retryConnDbMaxTimes; i++ {
		// log.Printf("getClient %s try --->: %d\n", db, i)
		if atomic.LoadUint32(&dbc.done) == 0 {
			// log.Printf("newClient try --->: %s\n", db)
			atomic.StoreUint32(&dbc.done, 1)
			if dbc.pointer, err = m.openMysql(dialect); err == nil && dbc.pointer != nil {
				m.setDefaultConfig(dbc.pointer)
			}
		}
		if err == nil {
			if dbc.pointer == nil {
				err = fmt.Errorf("mysql NewClient(%s) is nil", dialect.Db)
			} else {
				if sqlDB, err2 := dbc.pointer.DB(); err2 != nil {
					err = err2
				} else {
					err = sqlDB.Ping()
				}
			}
		}
		if err != nil {
			atomic.StoreUint32(&dbc.done, 0)
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Millisecond * time.Duration(i*200))
				continue
			} else {
				go func() {
					msg := fmt.Sprintf("mysql NewClient(%s) error: %v", dialect.Db, err)
					app.Logger.Error(msg)
					tool.PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s", msg), true)
				}()
				return nil
			}
		}
		break
	}
	return dbc.pointer
}

func (m *Model) SetMaxOpenConns(db *gorm.DB, num int) {
	if num > 0 && num <= maxOpenConns {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.SetMaxOpenConns(num)
		}
	}
}

func (m *Model) openMysql(dialect *app.MysqlDialect) (*gorm.DB, error) {
	if dialect == nil {
		return nil, errors.New("mysql dialect is nil")
	}
	charset := defaultCharset
	if dialect.Charset != "" {
		charset = dialect.Charset
	}
	parseTime := defaultParseTime
	if dialect.ParseTime != "" {
		parseTime = dialect.ParseTime
	}
	loc := defaultLoc
	if dialect.Loc != "" {
		loc = dialect.Loc
	}
	maxAllowedPacket := defaultMaxAllowedPacket
	if dialect.MaxAllowedPacket != "" {
		maxAllowedPacket = dialect.MaxAllowedPacket
	}
	timeout := defaultTimeout
	if dialect.Timeout != "" {
		timeout = dialect.Timeout
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&maxAllowedPacket=%s&timeout=%s",
		dialect.User, dialect.Pwd, dialect.Host, dialect.Port, dialect.Db,
		charset, parseTime, loc, maxAllowedPacket, timeout)
	return gorm.Open(mysql.Open(dns), &gorm.Config{})
}

func (m *Model) setDefaultConfig(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime 设置了连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
}
