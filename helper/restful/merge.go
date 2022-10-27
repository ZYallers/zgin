package restful

import (
	"github.com/ZYallers/zgin/types"
	"reflect"
	"strings"
	"sync"
)

var mtx sync.Mutex

func Merge(in ...types.Restful) types.Restful {
	mtx.Lock()
	defer mtx.Unlock()
	res := types.Restful{}
	for _, restful := range in {
		for path, handlers := range restful {
			if _, ok := res[path]; ok {
				continue
			}
			for i, hl := range handlers {
				if hl.Http == "" || hl.Handler == nil || hl.Method == "" {
					continue
				}
				if _, ok := reflect.ValueOf(hl.Handler).Type().MethodByName(hl.Method); !ok {
					continue
				}
				httpSlice := strings.Split(hl.Http, ",")
				httpMap := make(map[string]byte, len(httpSlice))
				for _, httpMethod := range httpSlice {
					httpMap[strings.ToUpper(httpMethod)] = 1
				}
				hl.Https = httpMap
				handlers[i] = hl
			}
			res[path] = handlers
		}
	}
	return res
}
