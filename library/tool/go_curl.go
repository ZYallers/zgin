package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 请求默认超时时间
const requestDefaultTimeout = 10 * time.Second

type goCurl struct {
	timeout        time.Duration
	requests       []interface{}
	reqCounter     int
	debug          bool
	Data           map[string]interface{}
	reqDepend      map[string]chan interface{}
	reqDependTimes map[string]uint64
	Result         map[string]interface{}
	resultChan     chan interface{}
	Runtime        time.Duration
}

func NewGoCurl() *goCurl {
	return &goCurl{
		timeout:        requestDefaultTimeout,
		reqDependTimes: make(map[string]uint64),
		resultChan:     make(chan interface{}),
		reqDepend:      make(map[string]chan interface{}),
		Result:         make(map[string]interface{}),
	}
}

func (gc *goCurl) Debug() *goCurl {
	gc.debug = true
	return gc
}

func (gc *goCurl) print(format string, v ...interface{}) *goCurl {
	if gc.debug {
		log.Printf(format, v...)
	}
	return gc
}

func (gc *goCurl) SetData(input string) *goCurl {
	// 不用json.Unmarshal方法，是因为解析后存在科学计数法问题
	d := json.NewDecoder(strings.NewReader(input))
	d.UseNumber()
	_ = d.Decode(&gc.Data)
	return gc
}

func (gc *goCurl) Done() (*goCurl, error) {
	var nowTime = time.Now()
	if val, ok := gc.Data["timeout"].(string); ok {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			gc.timeout = time.Duration(num) * time.Second
		}
	}

	if val, ok := gc.Data["requests"].([]interface{}); ok {
		gc.requests = val
	} else {
		return gc, errors.New(`lack of necessary parameters "requests"`)
	}
	gc.print("request length: %d, timeout: %v\n", len(gc.requests), gc.timeout)

	for _, req := range gc.requests {
		if params, ok := req.(map[string]interface{})["params"].(map[string]interface{}); ok {
			for _, param := range params {
				if param, ok := param.(map[string]interface{}); ok {
					dependId := param["depend_id"].(string)
					gc.reqDependTimes[dependId]++
				}
			}
		}
	}
	gc.print("reqDependTimes: %v\n", gc.reqDependTimes)

	for _, val := range gc.requests {
		if req, ok := val.(map[string]interface{}); ok {
			id := req["id"].(string)
			if dependTimes, ok := gc.reqDependTimes[id]; ok && dependTimes > 0 {
				gc.reqDepend[id] = make(chan interface{})
			}
			go gc.handler(req)
		}
	}

LOOP:
	for {
		select {
		case resp, ok := <-gc.resultChan:
			if ok {
				gc.reqCounter++
				gc.print("resultChan: %v, counter: %d\n", resp, gc.reqCounter)
				if ele, ok := resp.(map[string]string); ok {
					gc.Result[ele["id"]] = map[string]string{"data": ele["data"], "runtime": ele["runtime"]}
				}
				if gc.reqCounter >= len(gc.requests) {
					break LOOP
				}
			}
		case <-time.After(gc.timeout):
			gc.print("timeout: %s\n", gc.timeout)
			break LOOP
		}
	}
	gc.safeCloseChan(gc.resultChan)
	gc.Runtime = time.Since(nowTime)
	return gc, nil
}

func (gc *goCurl) handler(req map[string]interface{}) {
	var (
		nowTime = time.Now()
		id      string
	)
	if val, ok := req["id"].(string); ok {
		id = val
	} else {
		err := `missing necessary parameters "id"`
		gc.print(err)
		gc.safeSendChan(gc.resultChan, map[string]string{"id": id, "data": err, "runtime": time.Since(nowTime).String()})
		return
	}

	var url string
	if val, ok := req["url"].(string); ok {
		url = val
	} else {
		err := `missing necessary parameters "url"`
		gc.print(err)
		gc.safeSendChan(gc.resultChan, map[string]string{"id": id, "data": err, "runtime": time.Since(nowTime).String()})
		return
	}

	var httpMethod = http.MethodGet
	if val, ok := req["type"].(string); ok {
		httpMethod = strings.ToUpper(val)
	}

	var timeout = gc.timeout - 1*time.Second // 比全局的timeout少1秒
	if val, ok := req["timeout"].(string); ok {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			if to := time.Duration(num) * time.Second; to < timeout { // DIY的时间不能大于全局的时间
				timeout = to
			}
		}
	}

	headers := make(map[string]string)
	switch httpMethod {
	case http.MethodPost:
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	default:
		headers["Content-Type"] = "application/json;charset=utf-8"
	}
	if val, ok := req["headers"].(map[string]interface{}); ok {
		for k, v := range val {
			headers[k] = v.(string)
		}
	}

	queries := make(map[string]string)
	postData := make(map[string]interface{})
	if params, ok := req["params"].(map[string]interface{}); ok {
		for key, value := range params {
			var transfer interface{}
			if depend, ok := value.(map[string]interface{}); ok {
				dependId := depend["depend_id"].(string)
				if dependChan, ok := gc.reqDepend[dependId]; ok {
				LOOP:
					for {
						select {
						case resp, ok := <-dependChan:
							if ok {
								gc.print("dependChan: %s, resp: %s\n", dependId, resp)
								if res := gjson.Get(resp.(string), depend["depend_param"].(string)); res.Exists() {
									transfer = res.Value()
								}
								break LOOP
							}
						case <-time.After(timeout):
							gc.print("dependChan: %s, timeout: %s\n", dependId, timeout)
							break LOOP
						}
					}
				}
			} else {
				transfer = value
			}
			if transfer == nil {
				transfer = ""
			}
			if httpMethod == http.MethodGet {
				queries[key] = fmt.Sprintf("%v", transfer)
			} else {
				postData[key] = transfer
			}
		}
	}

	curl := NewRequest(url).SetMethod(httpMethod).SetTimeOut(timeout).SetHeaders(headers).SetQueries(queries).SetPostData(postData)
	gc.print("------>begin id: %s, url: %s, type: %s, headers: %#v, queries: %#v, postData: %#v\n", id, url, httpMethod, headers, queries, postData)

	var respBody string
	if resp, err := curl.Send(); err == nil {
		respBody = resp.Body
	} else {
		respBody = err.Error()
	}
	runtime := time.Since(nowTime).String()
	gc.print("<------end id: %s, runtime: %s, respBody: %v.\n", id, runtime, respBody)
	gc.safeSendChan(gc.resultChan, map[string]string{"id": id, "data": respBody, "runtime": runtime})

	if dependChan, ok := gc.reqDepend[id]; ok {
		if dependTimes, ok := gc.reqDependTimes[id]; ok && dependTimes > 0 {
			var i uint64
			for i = 0; i < dependTimes; i++ {
				gc.safeSendChan(dependChan, respBody)
			}
			gc.safeCloseChan(dependChan)
		}
	}
}

func (gc *goCurl) safeSendChan(ch chan<- interface{}, value interface{}) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}

func (gc *goCurl) safeCloseChan(ch chan interface{}) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = false
		}
	}()
	close(ch)
	return true
}
