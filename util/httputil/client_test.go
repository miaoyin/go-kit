package httputil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostJson(t *testing.T) {
	var result map[string]any
	clientUrl := "http://127.0.0.1:9111"
	rsp := JsonPost(
		clientUrl,
		map[string]any{
			"Name": "TestEvent",
			"Value": nil,
		},
	).CheckOK().UnmarshalBody(&result, json.Unmarshal)
	assert.Nil(t, rsp.Error())
	fmt.Println(result)
}


func TestPostBytes(t *testing.T) {
	var result []byte
	clientUrl := "http://127.0.0.1:9111"
	rsp := JsonPost(
		clientUrl,
		map[string]any{
			"Name": "TestEvent",
			"Value": nil,
		},
	)
	assert.Nil(t, rsp.CheckOK().UnmarshalBody(&result, json.Unmarshal).Error())
	fmt.Println(result)
}


func TestGet(t *testing.T) {
	clientUrl := "http://127.0.0.1:9111"
	options := make(url.Values)
	options.Add("a", "123")
	rsp := Get(clientUrl + "?" + options.Encode())
	data, err := rsp.ReadBody()
	assert.Nil(t, err)
	fmt.Println(rsp.StatusCode)
	fmt.Println(string(data))
}