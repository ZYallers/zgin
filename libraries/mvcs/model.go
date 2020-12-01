package mvcs

import (
	"bytes"
	"errors"
	"fmt"
	app "github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync/atomic"
	"time"
)

const (
	dialectConfig       = "?charset=utf8&loc=PRC&parseTime=true&maxAllowedPacket=0&timeout=10s"
	retryConnDbMaxTimes = 3
	maxOpenConns        = 1000
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
				err = dbc.pointer.DB().Ping()
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
		db.DB().SetMaxOpenConns(num)
	}
}

func (m *Model) openMysql(dialect *app.MysqlDialect) (*gorm.DB, error) {
	if dialect == nil {
		return nil, errors.New("mysql dialect is nil")
	}
	var tcp bytes.Buffer
	tcp.WriteString(dialect.User)
	tcp.WriteString(":")
	tcp.WriteString(dialect.Pwd)
	tcp.WriteString("@tcp(")
	tcp.WriteString(dialect.Host)
	tcp.WriteString(":")
	tcp.WriteString(dialect.Port)
	tcp.WriteString(")/")
	tcp.WriteString(dialect.Db)
	tcp.WriteString(dialectConfig)
	return gorm.Open("mysql", tcp.String())
}

func (m *Model) setDefaultConfig(db *gorm.DB) {
	db.DB().SetMaxOpenConns(8)
	db.DB().SetMaxIdleConns(2)
	db.DB().SetConnMaxLifetime(time.Second * 30)
}
