package types

import (
	"github.com/ZYallers/golib/types"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/golib/utils/mysql"
	"github.com/ZYallers/zgin/helper/dingtalk"
	"gorm.io/gorm"
)

type Model struct {
	mysql.Model
}

func (m *Model) New(dbc *types.DBCollector, dialect *types.MysqlDialect) *gorm.DB {
	db, err := m.NewMysql(dbc, dialect)
	if err != nil {
		logger.Use("mysql").Error(err.Error())
		dingtalk.PushSimpleMessage(err.Error(), true)
	}
	return db
}
