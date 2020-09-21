package mvcs

import (
	"errors"
	"fmt"
	app "github.com/ZYallers/go-frame/gin/application"
	"github.com/ZYallers/go-frame/gin/library/tool"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

const (
	retryConnRdsMaxTimes = 3
)

type Redis struct {
}

var test2 RdsCollector

type RdsCollector struct {
	once    sync.Once
	pointer *redis.Client
}

func (r *Redis) NewClient(rdc *RdsCollector, client *app.RedisClient) *redis.Client {
	var (
		err   error
		fatal bool
	)
	for i := 1; i <= retryConnRdsMaxTimes; i++ {
		//log.Printf("getClient %s try --->: %d\n", client, i)
		rdc.once.Do(func() {
			//log.Printf("newClient try --->: %s\n", client)
			rdc.pointer, err = r.newClient(client)
		})
		if err != nil {
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				rdc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		if err = rdc.pointer.Ping().Err(); err != nil {
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				rdc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		break
	}
	if fatal {
		defer func() {
			errMsg := fmt.Sprintf("get redis client of %s occur error: %s", client, err.Error())
			app.Logger.Error(errMsg)
			tool.PushSimpleMessage(errMsg, true)
		}()
		return nil
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

func (r *Redis) GetCache() *redis.Client {
	return r.NewClient(&test2, app.TestRedis)
}
