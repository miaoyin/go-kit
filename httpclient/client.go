package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

//----------------------------------------
//   选项
//----------------------------------------

type Option func(c *Client)

// WithTimeout 控制请求的超时
// 上下文context控制的是整个请求, context与timeout可共存
// timeout < context中的timeout
func WithTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

func WithHeader(k, v string) Option {
	return func(c *Client) {
		c.headers[k] = v
	}
}

func WithTransport(tr *http.Transport) Option {
	return func(c *Client) {
		c.client.Transport = tr
	}
}

//----------------------------------------
//   client
//----------------------------------------

var _ HTTPClient = (*Client)(nil)

// Client 客户端
//
//	声明为全局对象, 减少New创建
type Client struct {
	client  *http.Client
	timeout time.Duration
	headers map[string]string
}

func New(opts ...Option) *Client {
	tr := &http.Transport{
		MaxIdleConns:        100,              //默认100, 全局最多保留多少空闲连接
		MaxIdleConnsPerHost: 20,               //默认2, 单个host最多保留多少连接, 高平发需要调高
		IdleConnTimeout:     30 * time.Second, //默认90s, 减少僵死链接. 空闲连接放多久关闭
	}

	c := &Client{
		client: &http.Client{
			Transport: tr,
		},
		headers: make(map[string]string),
		timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}
	c.client.Timeout = c.timeout
	return c
}

func (c *Client) Do(ctx context.Context, method, reqURL string, body io.Reader, headers map[string]string) (*Response, error) {
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("new request err: %w", err)
	}
	//绑定ctx
	req = req.WithContext(ctx)

	// 公共header
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	// 单次请求临时header
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do err: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body err: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}, nil
}

func (c *Client) Get(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperGet(ctx, c, reqURL, query, headers)
}

func (c *Client) PostJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPostJSON(ctx, c, reqURL, reqBody, headers)
}

func (c *Client) PostForm(ctx context.Context, reqURL string, form map[string]string, headers map[string]string) (*Response, error) {
	return HelperPostForm(ctx, c, reqURL, form, headers)
}

func (c *Client) PutJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPutJSON(ctx, c, reqURL, reqBody, headers)
}

func (c *Client) Delete(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperDelete(ctx, c, reqURL, query, headers)
}

func (c *Client) Head(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperHead(ctx, c, reqURL, query, headers)
}

func (c *Client) PatchJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPatchJSON(ctx, c, reqURL, reqBody, headers)
}
