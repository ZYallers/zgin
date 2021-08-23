package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

// NewRequest
func NewRequest(url string) *Request {
	return &Request{Url: url, client: HttpClient, Timeout: DefaultHttpClientTimeout}
}

// SetMethod
func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

// SetUrl
func (r *Request) SetUrl(url string) *Request {
	r.Url = url
	return r
}

// SetHeaders
func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.Headers = headers
	return r
}

// 将用户自定义请求头添加到http.Request实例上
func (r *Request) setHeaders() *Request {
	for k, v := range r.Headers {
		r.request.Header.Set(k, v)
	}
	return r
}

// SetCookies
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

// SetPostData
func (r *Request) SetPostData(postData map[string]interface{}) *Request {
	if postData != nil {
		r.PostData = postData
		r.Body = nil
	}
	return r
}

// 设置post请求的提交数据
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

// SetBody
func (r *Request) SetBody(body io.Reader) *Request {
	if body != nil {
		r.Body = body
		r.PostData = nil
	}
	return r
}

// SetTimeOut
func (r *Request) SetTimeOut(timeout time.Duration) *Request {
	if timeout > 0 && timeout < DefaultHttpClientTimeout {
		r.Timeout = timeout
	}
	return r
}

// 发起get请求
func (r *Request) Get() (*Response, error) {
	return r.SetMethod(http.MethodGet).Send()
}

// 发起post请求
func (r *Request) Post() (*Response, error) {
	return r.SetMethod(http.MethodPost).Send()
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
	if bts, err := IoCopy(r.Raw.Body); err != nil {
		return err
	} else {
		r.Body = string(bts)
		return nil
	}
}
