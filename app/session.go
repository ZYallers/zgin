package app

import (
	"github.com/go-redis/redis"
	"time"
)

type SessionConfig struct {
	Client         func() *redis.Client
	TokenKey       string
	DataKey        string
	UpdateDuration time.Duration
	Expiration     int64
}

var Session = &SessionConfig{
	TokenKey:       "sess_token",
	DataKey:        "gin-gonic/gin/sessdata",
	UpdateDuration: 30 * time.Minute, // 30minutes
	Expiration:     6 * 30 * 86400,   // 6months
}
