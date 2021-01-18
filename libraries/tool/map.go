package tool

import (
	"reflect"
	"sort"
	"strings"
)

// Map / Slice 深度copy
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}
		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}
		return newSlice
	}
	return value
}

//经典排序返回a=1&b=1
func SortMapByKey(mp map[string]interface{}) string {
	if len(mp) == 0 {
		return ""
	}
	var newMp = make([]string, 0)
	for k := range mp {
		newMp = append(newMp, k)
	}
	sort.Strings(newMp)
	str := ""
	for _, v := range newMp {
		str += v + "=" + mp[v].(string) + "&"
	}
	return strings.TrimRight(str, "&")
}

//结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
