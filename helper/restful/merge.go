package restful

import (
	"fmt"
	"github.com/ZYallers/zgin/types"
	"reflect"
	"strings"
	"sync"
)

var mu sync.Mutex

func Merge(in ...types.Restful) types.Restful {
	mu.Lock()
	defer mu.Unlock()
	res := types.Restful{}
	for _, restful := range in {
		for path, restHandlers := range restful {
			if val, exist := res[path]; exist {
				panic(fmt.Errorf("restful path \"%s\" already exists: %+v", path, val))
			}
			for i, rh := range restHandlers {
				if rh.Http == "" || rh.Handler == nil || rh.Method == "" {
					panic(fmt.Errorf("restHandler attribute assignment is invalid: %+v", rh))
				}
				hmdSplit := strings.Split(rh.Http, ",")
				https := make(map[string]byte, len(hmdSplit))
				for _, httpMethod := range hmdSplit {
					https[strings.ToUpper(httpMethod)] = 1
				}
				rh.Https = https
				handlerValueOf := reflect.ValueOf(rh.Handler)
				if _, exist := handlerValueOf.Type().MethodByName(rh.Method); !exist {
					panic(fmt.Errorf("restHandler.Method does not exist: %+v\n", rh))
				}
				restHandlers[i] = rh
			}
			res[path] = restHandlers
		}
	}
	return res
}
