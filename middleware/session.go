package middleware

import (
	"fmt"
	"github.com/ZYallers/golib/funcs/php"
	"github.com/ZYallers/zgin/types"
)

func GetSessionData(ss *types.Session, token string) map[string]interface{} {
	if token == "" || ss.GetClientFunc == nil {
		return nil
	}
	client := ss.GetClientFunc()
	if client == nil {
		return nil
	}
	if str, _ := client.Get(ss.KeyPrefix + token).Result(); str != "" {
		fmt.Println("str:", str)
		return php.Unserialize(str)
	}
	return nil
}
