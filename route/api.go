package route

import (
	"fmt"
	"github.com/ZYallers/zgin/libraries/tool"
	"reflect"
	"strings"
	"sync"
)

var mu sync.Mutex

func Merge(in ...Restful) Restful {
	mu.Lock()
	defer mu.Unlock()

	res := Restful{}
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
				hmd := make(map[string]byte, len(hmdSplit))
				for _, httpMethod := range hmdSplit {
					hmd[strings.ToUpper(httpMethod)] = 1
				}
				rh.SetHttpMethod(hmd)

				methodName := tool.StrFirstToUpper(rh.Method)
				if _, exist := reflect.TypeOf(rh.Handler).MethodByName(methodName); !exist {
					panic(fmt.Errorf("restHandler.Method does not exist: %+v\n", rh))
				} else {
					rh.SetMethod(reflect.ValueOf(rh.Handler).MethodByName(methodName))
				}
				restHandlers[i] = rh
			}
			res[path] = restHandlers
		}
	}
	return res
}
