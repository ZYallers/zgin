package app

import (
	"github.com/go-redis/redis"
	"time"
)

type SessionConfig struct {
	Client         *redis.Client
	TokenKey       string
	DataKey        string
	UpdateDuration time.Duration
	Expiration     int64
}

var Session = &SessionConfig{
	TokenKey:       "gin-gonic/gin/sesstoken",
	DataKey:        "gin-gonic/gin/sessdata",
	UpdateDuration: 5 * time.Minute, // 5minutes
	Expiration:     6 * 30 * 86400,  // 6months
}
