package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BuildQueryURL 拼接query参数
func BuildQueryURL(rawURL string, query map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// JsonBody 请求json
func JsonBody(reqBody any, headers map[string]string) (io.Reader, error) {
	bs, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json;charset=utf-8"
	return bytes.NewBuffer(bs), nil
}

// FormBody 表单body
func FormBody(form map[string]string, headers map[string]string) (io.Reader, error) {
	val := url.Values{}
	for k, v := range form {
		val.Set(k, v)
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return bytes.NewBufferString(val.Encode()), nil
}

func HelperGet(ctx context.Context, d Doer, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	finalURL, err := BuildQueryURL(reqURL, query)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodGet, finalURL, nil, headers)
}

func HelperHead(ctx context.Context, d Doer, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	finalURL, err := BuildQueryURL(reqURL, query)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodHead, finalURL, nil, headers)
}

func HelperDelete(ctx context.Context, d Doer, reqURL string, query map[string]string, headers map[string]string) (*Response, error) {
	finalURL, err := BuildQueryURL(reqURL, query)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodDelete, finalURL, nil, headers)
}

func HelperPostJSON(ctx context.Context, d Doer, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	bodyR, err := JsonBody(reqBody, headers)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodPost, reqURL, bodyR, headers)
}

func HelperPutJSON(ctx context.Context, d Doer, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	bodyR, err := JsonBody(reqBody, headers)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodPut, reqURL, bodyR, headers)
}

func HelperPatchJSON(ctx context.Context, d Doer, reqURL string, reqBody any, headers map[string]string) (*Response, error) {
	bodyR, err := JsonBody(reqBody, headers)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodPatch, reqURL, bodyR, headers)
}

func HelperPostForm(ctx context.Context, d Doer, reqURL string, form map[string]string, headers map[string]string) (*Response, error) {
	bodyR, err := FormBody(form, headers)
	if err != nil {
		return nil, err
	}
	return d.Do(ctx, http.MethodPost, reqURL, bodyR, headers)
}
