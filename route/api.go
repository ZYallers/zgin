package route

import (
	"fmt"
	"github.com/ZYallers/zgin/libraries/tool"
	"reflect"
	"strings"
	"sync"
)

var (
	Api *Restful
	mu  sync.Mutex
)

func Merge(restful Restful) {
	mu.Lock()
	defer mu.Unlock()
	api := Restful{}
	if Api != nil {
		api = *Api
	}
	for k, restHandlers := range restful {
		for i, rh := range restHandlers {
			if rh.Http == "" || rh.Handler == nil || rh.Method == "" {
				panic(fmt.Errorf("restHandler attribute assignment is invalid: %+v", rh))
			}
			hmdSplit := strings.Split(rh.Http, ",")
			hmd := make(map[string]byte, len(hmdSplit))
			for _, httpMethod := range hmdSplit {
				hmd[strings.ToUpper(httpMethod)] = 1
			}
			rh.SetHttpMethod(hmd)

			typeOf := reflect.TypeOf(rh.Handler)
			if m, exist := typeOf.MethodByName(tool.StrFirstToUpper(rh.Method)); !exist {
				panic(fmt.Errorf("restHandler.Method does not exist: %+v\n", rh))
			} else {
				rh.SetMethod(m.Func)
			}
			restHandlers[i] = rh
		}
		api[k] = restHandlers
	}
	Api = &api
}
