package mvcs

import (
	"errors"
	"fmt"
	app "github.com/ZYallers/zgin/application"
	"github.com/ZYallers/zgin/library/tool"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

const (
	Dialect             = "charset=utf8&loc=PRC&parseTime=true&maxAllowedPacket=0&timeout=10s"
	retryConnDbMaxTimes = 3
)

type Model struct {
}

type DbCollector struct {
	once    sync.Once
	pointer *gorm.DB
}

var test DbCollector

func (m *Model) NewClient(dbc *DbCollector, dialect *app.MysqlDialect) *gorm.DB {
	var (
		err   error
		fatal bool
	)
	for i := 1; i <= retryConnDbMaxTimes; i++ {
		//log.Printf("getClient %s try --->: %d\n", db, i)
		dbc.once.Do(func() {
			//log.Printf("openMysql try --->: %s\n", db)
			if dbc.pointer, err = m.openMysql(dialect); err == nil {
				m.setDefaultConfig(dbc.pointer)
			}
		})
		if err != nil {
			if i < retryConnDbMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				dbc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		if err = dbc.pointer.DB().Ping(); err != nil {
			if i < retryConnDbMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				dbc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		break
	}
	if fatal {
		defer func() {
			errMsg := fmt.Sprintf("new client of mysql occur error: %s", err.Error())
			app.Logger.Error(errMsg)
			tool.PushSimpleMessage(errMsg, true)
		}()
		return nil
	}
	return dbc.pointer
}

func (m *Model) SetMaxOpenConns(db *gorm.DB, num int) {
	if num > 0 && num <= 1000 {
		db.DB().SetMaxOpenConns(num)
	}
}

func (m *Model) openMysql(dialect *app.MysqlDialect) (*gorm.DB, error) {
	if dialect == nil {
		return nil, errors.New("mysql dialect is nil")
	}
	tcp := dialect.User + ":" + dialect.Pwd + "@tcp(" + dialect.Host + ":" + dialect.Port + ")/" + dialect.Db + "?" + Dialect
	return gorm.Open("mysql", tcp)
}

func (m *Model) setDefaultConfig(db *gorm.DB) {
	db.DB().SetMaxOpenConns(8)
	db.DB().SetMaxIdleConns(2)
	db.DB().SetConnMaxLifetime(time.Second * 30)
}

func (m *Model) GetTest() *gorm.DB {
	return m.NewClient(&test, app.TestMysql)
}
