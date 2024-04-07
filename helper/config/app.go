package config

import (
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	singleton sync.Once
	cfg       map[string]interface{}
)

func ReadFile(args ...string) error {
	relativePath, configName, configType := ".", "cfg", "json"
	argsLen := len(args)
	if argsLen > 0 {
		relativePath = args[0]
	}
	if argsLen > 1 {
		configName = args[1]
	}
	if argsLen > 2 {
		configType = args[2]
	}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	_, filePath, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filePath), relativePath)
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	return viper.ReadInConfig()
}

func AllSettings() map[string]interface{} {
	singleton.Do(func() { cfg = viper.AllSettings() })
	return cfg
}

func Value(key string) interface{} {
	allCfg := AllSettings()
	if val, exist := allCfg[strings.ToLower(key)]; exist {
		return val
	}
	return nil
}

func AppMap() map[string]interface{} {
	if m := Value("app"); m != nil {
		return m.(map[string]interface{})
	}
	return nil
}

func AppValue(key string) interface{} {
	if am := AppMap(); am == nil {
		return nil
	} else if val, exist := am[key]; exist {
		return val
	}
	return nil
}
