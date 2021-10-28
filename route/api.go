package route

import (
	"fmt"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"reflect"
	"sort"
	"strconv"
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
				rh.http = hmd
				handlerValue := reflect.ValueOf(rh.Handler)
				if _, exist := handlerValue.Type().MethodByName(rh.Method); !exist {
					panic(fmt.Errorf("restHandler.Method does not exist: %+v\n", rh))
				} else {
					rh.method = handlerValue.MethodByName(rh.Method)
				}
				restHandlers[i] = rh
			}
			res[path] = restHandlers
		}
	}
	return res
}

func Register(des Restful, controllers ...mvcs.IController) Restful {
	mu.Lock()
	defer mu.Unlock()
	res := Restful{}
	for path, restHandlers := range des {
		res[path] = restHandlers
	}
	for _, controller := range controllers {
		contValue := reflect.ValueOf(controller)
		contName := contValue.Elem().Type().Name()
		if _, exist := contValue.Elem().Type().FieldByName("tag"); !exist {
			continue
		}
		tagVal := contValue.Elem().FieldByName("tag")
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
				panic(fmt.Errorf("restHandler.Path is empty: %s.%s\n", contName, methodName))
			}
			htt := fieldTagVal.Get("http")
			if htt == "" {
				panic(fmt.Errorf("restHandler.Http is empty: %s.%s\n", contName, methodName))
			}
			if _, exist := contValue.Type().MethodByName(methodName); !exist {
				panic(fmt.Errorf("restHandler.Method does not exist: %s.%s\n", contName, methodName))
			}
			httSplit := strings.Split(htt, ",")
			httpMap := make(map[string]byte, len(httSplit))
			for _, httpMethod := range httSplit {
				httpMap[strings.ToUpper(httpMethod)] = 1
			}

			ptr := reflect.New(contValue.Type().Elem())
			ptr.Elem().Set(contValue.Elem())
			newController := ptr.Interface().(mvcs.IController)

			resHandler := RestHandler{
				Path:    path,
				Http:    htt,
				Handler: newController,
				Version: fieldTagVal.Get("ver"),
				Signed:  fieldTagVal.Get("sign") == "on",
				Logged:  fieldTagVal.Get("login") == "on",
				http:    httpMap,
				method:  reflect.ValueOf(newController).MethodByName(methodName),
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
