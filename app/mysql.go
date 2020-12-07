package app

type MysqlDialect struct {
	User, Pwd, Host, Port, Db                          string
	Charset, Loc, ParseTime, MaxAllowedPacket, Timeout string
}
