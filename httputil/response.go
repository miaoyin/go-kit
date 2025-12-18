package httputil

import (
	"fmt"
	"io"
	"net/http"
)


func CreateResponse(requestFunc RequestFunc) *Response{
	return &Response{
		requestFunc: requestFunc,
	}
}

func ErrorResponse(err error) *Response{
	return &Response{
		err: err,
	}
}

type Response struct{
	*http.Response
	// requestFunc 请求函数
	requestFunc    RequestFunc
	err error
}

func (r *Response) Error() error{
	return r.err
}

func (r *Response) DoRequest() *Response {
	if r.err != nil {
		return r
	}
	if r.Response ==nil{
		r.Response, r.err = r.requestFunc()
	}
	return r
}

// ReadBody 读取body
func (r *Response) ReadBody() ([]byte, error) {
	r.DoRequest()
	if r.err!=nil{
		return nil, r.err
	}

	// 读取body
	rawData, err := io.ReadAll(r.Response.Body)
	defer r.Response.Body.Close()
	if err!=nil{
		return nil, fmt.Errorf("ReadBody %v", err)
	}
	return rawData, err
}

// UnmarshalBody 解析结果
func (r *Response) UnmarshalBody(result any, unmarshaller Unmarshaller) *Response {
	// 读取
	dataBytes, err := r.ReadBody()
	if err != nil {
		r.err = err
		return r
	}
	//解析
	if unmarshaller!=nil{
		if err = unmarshaller(dataBytes, result);err!=nil{
			r.err = fmt.Errorf("UnmarshalError error=%v, body=%s", err, string(dataBytes))
			return r
		}
	}
	return r
}

// CheckStatusCode 校验code
func (r *Response) CheckStatusCode(code int) *Response{
	r.DoRequest()
	if r.err!=nil{
		return r
	}
	if r.Response.StatusCode != code {
		r.err = fmt.Errorf("status_code=%d", r.Response.StatusCode)
	}
	return r
}

func (r *Response) CheckOK() *Response{
	return r.CheckStatusCode(http.StatusOK)
}


