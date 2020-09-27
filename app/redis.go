package app

type RedisClient struct {
	Host, Port, Pwd string
	Db              int
}
