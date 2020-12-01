package mvcs

import (
	"errors"
	"fmt"
	"github.com/ZYallers/zgin/app"
	"github.com/ZYallers/zgin/libraries/tool"
	"github.com/go-redis/redis"
	"sync/atomic"
	"time"
)

const (
	retryConnRdsMaxTimes = 3
)

type Redis struct {
}

type RdsCollector struct {
	done    uint32
	pointer *redis.Client
}

func (r *Redis) NewClient(rdc *RdsCollector, client *app.RedisClient) *redis.Client {
	defer tool.SafeDefer()
	var err error
	for i := 1; i <= retryConnRdsMaxTimes; i++ {
		// log.Printf("getClient %s try --->: %d\n", client, i)
		if atomic.LoadUint32(&rdc.done) == 0 {
			// log.Printf("newClient try --->: %s\n", client)
			atomic.StoreUint32(&rdc.done, 1)
			rdc.pointer, err = r.newClient(client)
		}
		if err == nil {
			if rdc.pointer == nil {
				err = fmt.Errorf("redis NewClient(%s:%s) is nil", client.Host, client.Port)
			} else {
				err = rdc.pointer.Ping().Err()
			}
		}
		if err != nil {
			atomic.StoreUint32(&rdc.done, 0)
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Millisecond * time.Duration(i*200))
				continue
			} else {
				go func() {
					msg := fmt.Sprintf("redis NewClient(%s:%s) error: %v", client.Host, client.Port, err)
					app.Logger.Error(msg)
					tool.PushSimpleMessage(fmt.Sprintf("recovery from panic:\n%s", msg), true)
				}()
				return nil
			}
		}
		break
	}
	return rdc.pointer
}

func (r *Redis) newClient(client *app.RedisClient) (*redis.Client, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	rds := redis.NewClient(&redis.Options{
		Addr:     client.Host + ":" + client.Port,
		Password: client.Pwd,
		DB:       client.Db,
	})
	return rds, nil
}
