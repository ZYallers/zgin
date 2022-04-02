package types

import "github.com/go-redis/redis"

type ICheck interface {
	GetVersion() (version, key string)
	GetSign() (secretKey, key, timeKey, dev string, expiration int64)
	GetSession() (fn func() *redis.Client, key, prefix string, expiration int64)
}
