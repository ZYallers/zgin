package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 默认超时时间
const defaultTimeout = 10 * time.Second

// Request构造类
type Request struct {
	client   *http.Client
	request  *http.Request
	Method   string
	Url      string
	Timeout  time.Duration
	Headers  map[string]string
	Cookies  map[string]string
	Queries  map[string]string
	PostData map[string]interface{}
	Body     io.Reader
}

// 创建一个Request实例
func NewRequest(url string) *Request {
	return &Request{Url: url, client: &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}, Timeout: defaultTimeout}
}

// 设置请求方法
func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

// 设置请求地址
func (r *Request) SetUrl(url string) *Request {
	r.Url = url
	return r
}

// 设置请求头
func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.Headers = headers
	return r
}

// 将用户自定义请求头添加到http.Request实例上
func (r *Request) setHeaders() *Request {
	var foundConnection, foundUserAgent bool
	for k, v := range r.Headers {
		r.request.Header.Set(k, v)
		switch k {
		case "Connection":
			foundConnection = true
		case "User-Agent":
			foundUserAgent = true
		}
	}
	if !foundConnection {
		r.request.Close = true
		r.request.Header.Set("Connection", "close")
	}
	if !foundUserAgent {
		r.request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	}
	return r
}

// 设置请求cookies
func (r *Request) SetCookies(cookies map[string]string) *Request {
	r.Cookies = cookies
	return r
}

// 将用户自定义cookies添加到http.Request实例上
func (r *Request) setCookies() *Request {
	for k, v := range r.Cookies {
		r.request.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return r
}

// 设置url查询参数
func (r *Request) SetQueries(queries map[string]string) *Request {
	r.Queries = queries
	return r
}

// 将用户自定义url查询参数添加到http.Request上
func (r *Request) setQueries() *Request {
	q := r.request.URL.Query()
	for k, v := range r.Queries {
		q.Add(k, v)
	}
	r.request.URL.RawQuery = q.Encode()
	return r
}

// 设置post请求的提交数据
func (r *Request) SetPostData(postData map[string]interface{}) *Request {
	if postData != nil {
		r.PostData = postData
		r.Body = nil
	}
	return r
}

// setPostData
func (r *Request) setPostData() (err error) {
	if ct, ok := r.Headers["Content-Type"]; ok {
		switch strings.ToLower(ct) {
		case "application/json", "application/json;charset=utf-8":
			var bts []byte
			if bts, err = json.Marshal(r.PostData); err == nil {
				r.Body = bytes.NewReader(bts)
			}
			return
		}
	}
	// 如果 Content-Type 不能匹配到，默认用 application/x-www-form-urlencoded 的方式处理
	postData := url.Values{}
	for k, v := range r.PostData {
		postData.Add(k, fmt.Sprintf("%v", v))
	}
	r.Body = strings.NewReader(postData.Encode())
	return
}

// 发起get请求
func (r *Request) Get() (*Response, error) {
	return r.SetMethod(http.MethodGet).Send()
}

// 发起post请求
func (r *Request) Post() (*Response, error) {
	return r.SetMethod(http.MethodPost).Send()
}

// SetDialTimeOut
func (r *Request) SetTimeOut(timeout time.Duration) *Request {
	if timeout > 0 && timeout < defaultTimeout {
		r.Timeout = timeout
	}
	return r
}

// SetBody
func (r *Request) SetBody(body io.Reader) *Request {
	if body != nil {
		r.Body = body
		r.PostData = nil
	}
	return r
}

// 发起请求
func (r *Request) Send() (*Response, error) {
	if r.PostData != nil {
		if err := r.setPostData(); err != nil {
			return nil, err
		}
	}
	if req, err := http.NewRequest(r.Method, r.Url, r.Body); err != nil {
		return nil, err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
		defer cancel()
		r.request = req.WithContext(ctx)
	}

	r.setHeaders().setCookies().setQueries()

	if resp, err := r.client.Do(r.request); err != nil {
		return nil, err
	} else {
		res := NewResponse()
		res.Raw = resp
		defer res.Raw.Body.Close()
		if err := res.parseBody(); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
}

// Response 构造类
type Response struct {
	Raw     *http.Response
	Headers map[string]string
	Body    string
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) StatusCode() int {
	if r.Raw == nil {
		return 0
	}
	return r.Raw.StatusCode
}

func (r *Response) IsOk() bool {
	return r.StatusCode() == http.StatusOK
}

func (r *Response) parseHeaders() {
	headers := map[string]string{}
	for k, v := range r.Raw.Header {
		headers[k] = v[0]
	}
	r.Headers = headers
}

func (r *Response) parseBody() error {
	if bts, err := ioutil.ReadAll(r.Raw.Body); err != nil {
		return err
	} else {
		r.Body = string(bts)
		return nil
	}
}
