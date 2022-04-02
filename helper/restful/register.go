package restful

import (
	"fmt"
	"github.com/ZYallers/zgin/types"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var mu2 sync.Mutex

func Register(des types.Restful, controllers ...types.IController) types.Restful {
	mu2.Lock()
	defer mu2.Unlock()
	res := types.Restful{}
	for path, restHandlers := range des {
		res[path] = restHandlers
	}
	for _, controller := range controllers {
		controllerValueOf := reflect.ValueOf(controller)
		controllerName := controllerValueOf.Elem().Type().Name()
		if _, exist := controllerValueOf.Elem().Type().FieldByName("tag"); !exist {
			continue
		}
		tagVal := controllerValueOf.Elem().FieldByName("tag")
		if tagVal.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < tagVal.NumField(); i++ {
			if tagVal.Field(i).Kind() != reflect.Func {
				continue
			}
			methodName := tagVal.Type().Field(i).Name
			fieldTagVal := tagVal.Type().Field(i).Tag
			path := fieldTagVal.Get("path")
			if path == "" {
				panic(fmt.Errorf("restHandler.Path is empty: %s.%s\n", controllerName, methodName))
			}
			htt := fieldTagVal.Get("http")
			if htt == "" {
				panic(fmt.Errorf("restHandler.Http is empty: %s.%s\n", controllerName, methodName))
			}
			if _, exist := controllerValueOf.Type().MethodByName(methodName); !exist {
				panic(fmt.Errorf("restHandler.Method does not exist: %s.%s\n", controllerName, methodName))
			}
			httSplit := strings.Split(htt, ",")
			httpMap := make(map[string]byte, len(httSplit))
			for _, httpMethod := range httSplit {
				httpMap[strings.ToUpper(httpMethod)] = 1
			}
			resHandler := types.RestHandler{
				Path:    path,
				Http:    htt,
				Https:   httpMap,
				Handler: controller,
				Method:  methodName,
				Version: fieldTagVal.Get("ver"),
				Signed:  fieldTagVal.Get("sign") == "on",
				Logged:  fieldTagVal.Get("login") == "on",
			}
			if sortStr := fieldTagVal.Get("sort"); sortStr != "" {
				if sortInt, err := strconv.Atoi(sortStr); err != nil {
					panic(fmt.Errorf("restHandler sort is invalid: %s", sortStr))
				} else {
					resHandler.Sort = sortInt
				}
			}
			res[path] = append(res[path], resHandler)
			if len(res[path]) > 1 {
				resHandlers := res[path]
				sort.Sort(resHandlers)
				res[path] = resHandlers
			}
		}
	}
	return res
}
