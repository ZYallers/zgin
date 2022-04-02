package types

import (
	"github.com/ZYallers/golib/types"
	"github.com/ZYallers/golib/utils/logger"
	"github.com/ZYallers/golib/utils/redis"
	"github.com/ZYallers/zgin/helper/dingtalk"
	redis2 "github.com/go-redis/redis"
)

type Redis struct {
	redis.Redis
}

func (r *Redis) New(rdc *types.RedisCollector, client *types.RedisClient) *redis2.Client {
	c, err := r.NewRedis(rdc, client)
	if err != nil {
		logger.Use("redis").Error(err.Error())
		dingtalk.PushSimpleMessage(err.Error(), true)
	}
	return c
}
