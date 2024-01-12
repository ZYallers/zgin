package restful

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ZYallers/zgin/types"
)

var rtx sync.Mutex

func Register(des types.Restful, controllers ...types.IController) types.Restful {
	rtx.Lock()
	defer rtx.Unlock()
	res := types.Restful{}
	for path, restHandlers := range des {
		res[path] = restHandlers
	}
	for _, controller := range controllers {
		controllerValue := reflect.ValueOf(controller)
		if _, ok := controllerValue.Elem().Type().FieldByName("tag"); !ok {
			continue
		}
		tagVal := controllerValue.Elem().FieldByName("tag")
		if tagVal.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < tagVal.NumField(); i++ {
			if tagVal.Field(i).Kind() != reflect.Func {
				continue
			}
			methodName := tagVal.Type().Field(i).Name
			tagValue := tagVal.Type().Field(i).Tag
			path := tagValue.Get("path")
			https := tagValue.Get("http")
			if path == "" || https == "" {
				continue
			}
			if _, ok := controllerValue.Type().MethodByName(methodName); !ok {
				continue
			}
			httpSlice := strings.Split(https, ",")
			httpMap := make(map[string]byte, len(httpSlice))
			for _, httpMethod := range httpSlice {
				httpMap[strings.ToUpper(httpMethod)] = 1
			}
			restVer := types.RestVersion{Value: tagValue.Get("ver")}
			if verLen := len(restVer.Value); verLen > 0 && restVer.Value[verLen-1:] == "+" {
				restVer.Plus = true
				restVer.Value = restVer.Value[0 : verLen-1]
			}
			resHandler := &types.RestHandler{
				Path:    path,
				Http:    https,
				Https:   httpMap,
				Handler: controller,
				Method:  methodName,
				Version: restVer,
				Sign:    tagValue.Get("sign") == "on",
				Login:   tagValue.Get("login") == "on",
			}
			if s := tagValue.Get("sort"); s != "" {
				if val, err := strconv.Atoi(s); err == nil {
					resHandler.Sort = val
				}
			}
			res[path] = append(res[path], resHandler)
			if len(res[path]) > 1 {
				sort.Slice(res[path], func(i, j int) bool {
					return res[path][i].Sort > res[path][j].Sort
				})
			}
		}
	}
	return res
}
