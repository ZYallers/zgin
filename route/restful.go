package route

import (
	"github.com/ZYallers/zgin/libraries/mvcs"
	"reflect"
)

type Restful map[string][]RestHandler

type RestHandler struct {
	Version    string           // Version
	Http       string           // HttpMethod
	Method     string           // Method
	Handler    mvcs.IController // Handler
	method     reflect.Value    // method
	httpMethod map[string]byte  // httpMethod
	Signed     bool             // 签名验证
	Logged     bool             // 登录验证
}

// SetHttpMethod
func (rh *RestHandler) SetHttpMethod(m map[string]byte) {
	rh.httpMethod = m
}

// GetHttpMethod
func (rh *RestHandler) GetHttpMethod() map[string]byte {
	return rh.httpMethod
}

// SetHttpMethod
func (rh *RestHandler) SetMethod(r reflect.Value) {
	rh.method = r
}

// GetMethod
func (rh *RestHandler) GetMethod() reflect.Value {
	return rh.method
}

// CallMethod
func (rh *RestHandler) CallMethod() {
	rh.method.Call(nil)
}
