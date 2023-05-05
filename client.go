package hhttp

import (
	"context"
	"net/http"
	"time"
)

type hiclient struct {
	client *http.Client
	opt    Options
}
type RetryErrorFunc func(ctx context.Context, r *Request) error

type Header map[string][]string

var defaultTrimChars = string([]byte{
	'\t', // Tab.
	'\v', // Vertical tab.
	'\n', // New line (line feed).
	'\r', // Carriage return.
	'\f', // New page.
	' ',  // Ordinary space.
	0x00, // NUL-byte.
	0x85, // Delete.
	0xA0, // Non-breaking space.
})
var client = hiclient{
	client: &http.Client{
		Timeout: 30 * time.Second,
	},
	opt: Options{
		retryCount: 0,
		retryWait:  time.Duration(0),
		retryError: func(ctx context.Context, r *Request) error {
			return nil
		},
		timeout: 5 * time.Second,
	},
}

// Load 设置client的全局参数
func Load(opts ...Option) {
	for _, o := range opts {
		o(&client.opt)
	}
}

// SetHeader 以k-v格式设置header
func (r *Request) SetHeader(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

// SetHeaders 设置header参数 Header
// 例如: c.Headers(header)
func (r *Request) SetHeaders(headers Header) *Request {
	for k, header := range headers {
		for _, he := range header {
			r.header.Add(k, he)
		}
	}
	return r
}

// SetCookies 设置cookie
func (r *Request) SetCookies(hc ...*http.Cookie) *Request {
	r.cookies = append(r.cookies, hc...)
	return r
}

// SetTimeout 修改http.Client的超时时间，该超时时间默认是30s，优先级大于通过context设置的超时时间
func (r *Request) SetTimeout(t time.Duration) *Request {
	r.client.client.Timeout = t
	return r
}
