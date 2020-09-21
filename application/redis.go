package app

import (
	"os"
)

type RedisClient struct {
	Host, Port, Pwd string
	Db              int
}

var (
	TestRedis = &RedisClient{
		Host: os.Getenv("redis_test_host"),
		Port: os.Getenv("redis_test_port"),
		Pwd:  os.Getenv("redis_test_password"),
		Db:   0,
	}
)
