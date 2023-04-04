package hhttp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	// API接口前缀
	baseUrl string
	method  string
	header  *http.Header
	cookies []*http.Cookie
	client  hiclient
	opt     Options
}

type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	// ProtoMajor      int    // e.g. 1
	// ProtoMinor      int    // e.g. 0
	Header Header
	Body   []byte
	// ContentLength   int64
	Close           bool
	RequestDuration time.Duration
}

// New 设置公共参数
func New(opts ...Option) *Request {
	request := Request{
		baseUrl: "",
		method:  "",
		header:  &http.Header{},
		cookies: []*http.Cookie{},
		client:  client,
	}
	opt := client.opt

	for _, o := range opts {
		o(&opt)
	}

	request.opt = opt
	return &request
}

type HiHTTP interface {
	Get(ctx context.Context, urlStr string, data ...Param) (*Response, error)
	Post(ctx context.Context, urlStr string, p Payload) (*Response, error)
	Put(ctx context.Context, urlStr string, p Payload) (*Response, error)
	Delete(ctx context.Context, urlStr string, data ...Param) (*Response, error)
	Patch(ctx context.Context, urlStr string, p Payload) (*Response, error)
}

func (r *Request) execute(ctx context.Context, payload io.Reader) (*Response, error) {
	r.baseUrl = strings.Trim(r.baseUrl, defaultTrimChars)

	start := time.Now() // 记录执行的开始时间
	httpCtx, cancel := context.WithTimeout(ctx, r.opt.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(httpCtx, r.method, r.baseUrl, payload)
	if err != nil {
		return nil, err
	}
	if r.header != nil {
		req.Header = *r.header
	}
	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	res, err := r.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// 设置header
	header := Header{}
	for k, v := range res.Header {
		header[k] = v
	}

	response := Response{
		Status:          res.Status,
		StatusCode:      res.StatusCode,
		Proto:           res.Proto,
		Header:          header,
		Close:           res.Close,
		RequestDuration: time.Since(start),
	}
	response.Body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get 发送一个Get请求
// 也可以把参数直接放到URL后面，则data不传即可
func (r *Request) Get(ctx context.Context, urlStr string, data ...Param) (*Response, error) {
	r.baseUrl = urlStr
	// 如果参数直接放在URL后面，则Param为空，不必拼接query参数
	if len(data) > 0 {
		r.baseUrl += "?" + mergeParams(data...)
	}
	r.method = GET
	req, err := r.execute(ctx, nil)
	if err != nil {
		r.retry(ctx, nil)
	}
	return req, err
}

// Post 发送一个POST请求
func (r *Request) Post(ctx context.Context, urlStr string, p Payload) (*Response, error) {
	r.baseUrl = urlStr
	r.method = POST
	if p == nil {
		p = &formPayload{}
	}
	if r.header.Get(SerializationType) == "" {
		r.header.Add(SerializationType, p.ContentType())
	}
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}

// Put 发送Put请求
func (r *Request) Put(ctx context.Context, urlStr string, p Payload) (*Response, error) {
	r.baseUrl = urlStr
	r.method = PUT
	if p == nil {
		p = &formPayload{}
	}
	if r.header.Get(SerializationType) == "" {
		r.header.Add(SerializationType, p.ContentType())
	}
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}

// Delete 发送一个delete请求
func (r *Request) Delete(ctx context.Context, urlStr string, data ...Param) (*Response, error) {
	r.baseUrl = urlStr
	if len(data) > 0 {
		r.baseUrl += "?" + mergeParams(data...)
	}
	r.method = DELETE
	req, err := r.execute(ctx, nil)
	if err != nil {
		r.retry(ctx, nil)
	}
	return req, err
}

// Patch 发送patch请求
func (r *Request) Patch(ctx context.Context, urlStr string, p Payload) (*Response, error) {
	r.baseUrl = urlStr
	r.method = PATCH
	if p == nil {
		p = &formPayload{}
	}
	if r.header.Get(SerializationType) == "" {
		r.header.Add(SerializationType, p.ContentType())
	}
	req, err := r.execute(ctx, p.Serialize())
	if err != nil {
		r.retry(ctx, p.Serialize())
	}
	return req, err
}

// BodyString 直接返回一个string格式的body
func (r *Response) BodyString() string {
	return string(r.Body)
}
