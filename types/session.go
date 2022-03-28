package types

import (
	"github.com/go-redis/redis"
)

type Session struct {
	TokenKey      string
	KeyPrefix     string
	Expiration    int64
	GetClientFunc func() *redis.Client
}

var DefaultSession = &Session{
	TokenKey:   "sess_token",
	KeyPrefix:  "ci_session:",
	Expiration: 6 * 30 * 86400, // 6 months
}
