package nats

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	natsUrl = "nats://user:passwd@127.0.0.1:4222"
	natsName = "nats"
	natsSubject = "subject"
)

func TestNewClientSuccess(t *testing.T)  {
	//测试成功
	testConfigs := []Configuration{
		{Enable: true, Url:natsUrl},
		{Enable: true, Url:""},
		{Enable: true, Url:"nats://127.0.0.1"},
		{Enable: true, Url:"xxx://127.0.0.1"},
	}
	for _, testConfig := range testConfigs {
		//创建启动
		client := NewClient(natsName, &testConfig)
		err := client.Start(context.Background())
		require.Truef(t, err == nil, "Error starting client: %v", err)

		//发消息
		err = client.Publish(natsSubject, []byte("demo"))
		require.Truef(t, err == nil, "PublishError: %v", err)
	}
}

func TestNewClientFail(t *testing.T)  {
	//测试报错
	testConfigs := []Configuration{
		{Enable: false, Url:natsUrl},
		{Enable: true, Url:"nats://127.0.0.1:4223"},
	}
	for _, testConfig := range testConfigs {
		//创建启动
		client := NewClient(natsName, &testConfig)
		err := client.Start(context.Background())
		require.Truef(t, err != nil, "starting client err==nil")

		//发消息
		err = client.Publish(natsSubject, []byte("demo"))
		require.Truef(t, err != nil, "Publish err==nil")
	}
}

//TestClientRestart 重启
func TestClientRestart(t *testing.T)  {
	client := NewClient(natsName, &Configuration{
		Enable: true,
		Url:natsUrl,
	})
	for i:=0;i<100;i++{
		//启动
		err := client.Start(nil)
		require.Truef(t, err == nil, "Start err!=nil")
		//发消息
		err = client.Publish(natsSubject, []byte("demo"))
		require.Truef(t, err == nil, "Publish err!=nil")
		//停止
		err = client.Close(nil)
		require.Truef(t, err == nil, "Close err!=nil")
		sleepTime := time.Millisecond*time.Duration(rand.Int31n(1000))
		fmt.Println(i, sleepTime)
		time.Sleep(sleepTime)
	}
}
