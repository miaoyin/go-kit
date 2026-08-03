package nats

import (
	"errors"
	"reflect"

	"github.com/miaoyin/go-kit/module"

	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var _ module.Module[Configer] = &Client{}


var (
	ErrNatsConnectionClosed = errors.New("nats connection closed")
)


var (
	HeaderMsgID       = "MsgId"
	HeaderTimestamp   = "Timestamp"
	HeaderSource      = "Source"
	HeaderSubject     = "Subject"
	HeaderMsgType     = "MsgType"
	HeaderVersion     = "DataVersion"
)

//Configuration 基础配置
type Configuration struct{
	Enable  bool
	Url     string
	Options []nats.Option  `json:"-" toml:"-"`
}

func (c Configuration) GetBase() Configuration{
	return c
}

type Configer interface{
	GetBase() Configuration
}

type Client struct{
	*module.BaseModule[Configer]
	subMgr *SubscriberManager
	conn    *nats.Conn
	mu sync.RWMutex
}

//NewClient 客户端
//  1. 旧连接延迟close, 不能直接close
//  2. View确保conn生命周期, 不要直接GetConn
func NewClient(name string, config Configer) *Client {
	//限定为struct
	rt := reflect.TypeOf(config)
	if rt.Kind() == reflect.Ptr {
		panic("config must be a struct")
	}
    return &Client{
		BaseModule: module.NewBaseModule(name, config),
		subMgr: NewSubscriberManager(),
    }
}


//
//Start 连接nats
// 用途
// 	1. 重复Start，会重连
// 注意
// 	1.rebuilding 减少并发
// 	2.rebuilding 快速返回，避免等锁
func (c *Client) Start(_ context.Context) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	//是否enable
	config, _ := c.GetConfig()
	baseConfig := (*config).GetBase()
	if !baseConfig.Enable {
		return module.ErrModuleDisabled
	}

	//状态校验
	if err := c.CheckStart();err!=nil{
		return err
	}

	c.mu.Unlock()
	newConn, err := nats.Connect(baseConfig.Url, baseConfig.Options...)
	c.mu.Lock()
	if err != nil {
		return err
	}
	//状态校验
	if err = c.CheckStart();err!=nil{
		return err
	}
	//替换连接
	c.setConn(newConn)
	c.DoStart()
	return nil
}

//replaceConn 替换连接
func (c *Client) setConn(newConn *nats.Conn){
	old := c.conn
	c.conn = newConn
	//重新订阅
	go func() {
		_ = c.subMgr.ResubscribeAll(newConn)
	}()

	// 异步关闭旧连接
	if old!=nil{
		_ = old.Flush()
		go func() {
			time.Sleep(2 * time.Second)
			old.Close()
		}()
	}
}

//SetConn running中替换连接
func (c *Client) SetConn(newConn *nats.Conn) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.IsRunning(){
		return module.ErrNotRunning
	}
	c.setConn(newConn)
	return nil
}

func (c *Client) Close(_ context.Context) error{
    c.mu.Lock()
    defer c.mu.Unlock()
	if err := c.CheckClose();err!=nil{
		return err
	}
    if c.conn!=nil{
		_ = c.conn.Flush()
		//取消订阅
		c.subMgr.CloseAll()
		//关闭连接
        c.conn.Close()
    }
	//支持重启
	c.DoInit()
    return nil
}

func (c *Client) View(fn func(conn *nats.Conn) error) error{
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.IsRunning(){
		return module.ErrNotRunning
	}
	if c.conn == nil || c.conn.IsClosed() {
		return ErrNatsConnectionClosed
	}
	return fn(c.conn)
}

//SubscriberManager 订阅管理
// 不支持提供订阅接口
func (c *Client) SubscriberManager() *SubscriberManager{
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.subMgr
}

//Publish 发布消息
//  1.快速失败优先, 不提供失败重试, 业务层处理
//  2.并发调用时,会产生大量重连
func (c *Client) Publish(subj string, data []byte) error {
    return c.View(func(conn *nats.Conn) error {
        return conn.Publish(subj, data)
    })
}

//PublishMsg 发布消息
//  使用者统计
func (c *Client) PublishMsg(msg *nats.Msg) error {
    return c.View(func(conn *nats.Conn) error {
        return conn.PublishMsg(msg)
    })
}

func (c *Client) Request(subj string, data []byte, timeout time.Duration) (msg *nats.Msg, err error){
    err = c.View(func(conn *nats.Conn) error {
        msg, err = conn.Request(subj, data, timeout)
        return err
    })
    return msg, err
}

func (c *Client) RequestMsg(reqMsg *nats.Msg, timeout time.Duration) (msg *nats.Msg, err error){
    err = c.View(func(conn *nats.Conn) error {
        msg, err = conn.RequestMsg(reqMsg, timeout)
        return err
    })
    return msg, err
}

//Flush 默认10s
func (c *Client) Flush() error{
    return c.View(func(conn *nats.Conn) error {
        return conn.Flush()
    })
}
