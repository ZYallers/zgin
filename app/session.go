package app

import (
	"github.com/go-redis/redis"
	"time"
)

type SessionConfig struct {
	TokenKey       string
	DataKey        string
	Expiration     int64
	UpdateDuration time.Duration
	GetClientFunc  func() *redis.Client
}

var Session = &SessionConfig{
	TokenKey:       "sess_token",
	DataKey:        "gin-gonic/gin/sessdata",
	UpdateDuration: 30 * time.Minute, // 30minutes
	Expiration:     6 * 30 * 86400,   // 6months
}
