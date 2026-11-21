package httpclient

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/avast/retry-go/v4"
)

var _ HTTPClient = (*wrapperClient)(nil)

// wrapperClient 扩展样例
type wrapperClient struct {
	next HTTPClient

	retryCfg *retryConfig
}

type retryConfig struct {
	attempts  uint
	delay     time.Duration
	retryable func(resp *Response, err error) bool
}

type WrapperOption func(w *wrapperClient)

func WithRetry(attempts uint, delay time.Duration) WrapperOption {
	return func(w *wrapperClient) {
		w.retryCfg = &retryConfig{
			attempts: attempts,
			delay:    delay,
			retryable: func(resp *Response, err error) bool {
				if err != nil {
					return true
				}
				if resp.StatusCode >= 500 {
					return true
				}
				return false
			},
		}
	}
}

func Wrap(next HTTPClient, opts ...WrapperOption) HTTPClient {
	w := &wrapperClient{next: next}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *wrapperClient) Do(ctx context.Context, method, reqURL string, body io.Reader, headers map[string]string) (*Response, error) {
	var resp *Response
	var err error
	if w.retryCfg != nil {
		err = retry.Do(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			resp, err = w.next.Do(ctx, method, reqURL, body, headers)
			if w.retryCfg.retryable(resp, err) {
				return fmt.Errorf("retryable")
			}
			return nil
		},
			retry.Attempts(w.retryCfg.attempts),
			retry.Delay(w.retryCfg.delay),
			retry.Context(ctx),
		)
	} else {
		resp, err = w.next.Do(ctx, method, reqURL, body, headers)
	}
	return resp, err
}

func (w *wrapperClient) Get(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperGet(ctx, w, reqURL, query, headers)
}

func (w *wrapperClient) PostJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPostJSON(ctx, w, reqURL, reqBody, headers)
}

func (w *wrapperClient) PostForm(ctx context.Context, reqURL string, form map[string]string, headers map[string]string) (*Response, error) {
	return HelperPostForm(ctx, w, reqURL, form, headers)
}

func (w *wrapperClient) PutJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPutJSON(ctx, w, reqURL, reqBody, headers)
}

func (w *wrapperClient) Delete(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperDelete(ctx, w, reqURL, query, headers)
}

func (w *wrapperClient) Head(ctx context.Context, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	return HelperHead(ctx, w, reqURL, query, headers)
}

func (w *wrapperClient) PatchJSON(ctx context.Context, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	return HelperPatchJSON(ctx, w, reqURL, reqBody, headers)
}
