package app

import (
	"os"
)

type MysqlDialect struct {
	User, Pwd, Host, Port, Db string
}

var (
	TestMysql = &MysqlDialect{
		Host: os.Getenv("mysql_test_host"),
		User: os.Getenv("mysql_test_username"),
		Pwd:  os.Getenv("mysql_test_password"),
		Db:   os.Getenv("mysql_test_database"),
		Port: os.Getenv("mysql_test_port"),
	}
)
