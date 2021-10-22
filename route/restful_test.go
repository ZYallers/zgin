// @Title restful_test.go
// @Description restful_test.go
// @Author zhongyongbiao 2021/10/20 下午12:06
// @Software GoLand
package route

import (
	v000 "github.com/ZYallers/zgin/controller/v000"
	v110 "github.com/ZYallers/zgin/controller/v110"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	var api Restful
	api = Register(nil, &v000.Index{}, &v110.Index{})
	for path, restHandlers := range api {
		log.Printf("path: %s\n", path)
		for _, resHandler := range restHandlers {
			log.Printf("%+v\n", resHandler)
		}
	}
}
