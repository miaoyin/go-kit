package redis

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/miaoyin/go-kit/module"
	"github.com/redis/go-redis/v9"
)


var _ module.Module[Configer] = &Client{}

var ErrClientNil = errors.New("redis client is nil")

type StandaloneConfiguration struct {
    // 127.0.0.1:6379
    Addr string
    DB   int
}

type SentinelConfiguration struct {
    //主节点名称
    MasterName string
    //哨兵节点地址列表
    Addrs      []string
    //哨兵密码
    SentinelPassword string
}

type ClusterConfiguration struct {
    //集群所有节点
    Addrs          []string
    //是否仅读从节点
    ReadOnly       bool
    //按延迟路由
    RouteByLatency bool
}

type Configuration struct{
    //是否启用
    Enable  bool
    // 模式标识，必填：standalone / sentinel / cluster
    Mode string `yaml:"mode" toml:"mode" json:"mode"`

    // 密码
    Password     string
    Username     string
    DialTimeout  time.Duration
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    // 连接池配置
    PoolSize     int
    MinIdleConns int
    MaxRetries   int

    //单机模式独有配置
    Standalone *StandaloneConfiguration
    //哨兵模式独有配置
    Sentinel  *SentinelConfiguration
    //集群模式独有配置
    Cluster  *ClusterConfiguration
}

func (c Configuration) GetBase() Configuration{
	return c
}

type Configer interface{
	GetBase() Configuration
}


//NewRedisClient 创建通用client
func NewRedisClient(cfg Configuration) redis.UniversalClient {
    switch cfg.Mode {
    case "sentinel":
        //哨兵模式
        sentinelOpts := &redis.FailoverOptions{
            SentinelAddrs: cfg.Sentinel.Addrs,
            SentinelPassword: cfg.Sentinel.SentinelPassword,
            MasterName:  cfg.Sentinel.MasterName,
            Username:     cfg.Username,
            Password:     cfg.Password,
            DialTimeout:  cfg.DialTimeout,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            PoolSize:     cfg.PoolSize,
            MinIdleConns: cfg.MinIdleConns,
            MaxRetries:   cfg.MaxRetries,
        }
        return redis.NewFailoverClient(sentinelOpts)
    case "cluster":
        //主从模式
        clusterOpts := &redis.ClusterOptions{
            Addrs:        cfg.Cluster.Addrs,
            ReadOnly:     cfg.Cluster.ReadOnly,
            RouteByLatency: cfg.Cluster.RouteByLatency,
            Username:     cfg.Username,
            Password:     cfg.Password,
            DialTimeout:  cfg.DialTimeout,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            PoolSize:     cfg.PoolSize,
            MinIdleConns: cfg.MinIdleConns,
            MaxRetries:   cfg.MaxRetries,
        }
        return redis.NewClusterClient(clusterOpts)
    default:
        //单机模式
        baseOpts :=&redis.Options{
            DB: cfg.Standalone.DB,
            Addr: cfg.Standalone.Addr,
            Username: cfg.Username,
            Password: cfg.Password,
            DialTimeout:  cfg.DialTimeout,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            PoolSize:     cfg.PoolSize,
            MinIdleConns: cfg.MinIdleConns,
            MaxRetries:   cfg.MaxRetries,
        }
        return redis.NewClient(baseOpts)
    }
}

//Client 连接
// 1.实现module接口
// 2.支持重启(redis连接可更新)
type Client struct{
	*module.BaseModule[Configer]
    rdb    redis.UniversalClient
    mu     sync.RWMutex
}

func NewClient(name string, config Configer)*Client{
	//限定为struct
	rt := reflect.TypeOf(config)
	if rt.Kind() == reflect.Ptr {
		panic("config must be a struct")
	}

    return &Client{
		BaseModule: module.NewBaseModule(name, config),
    }
}

func (c *Client) Start(_ context.Context) error{
    c.mu.Lock()
    defer c.mu.Unlock()
	config, _ := c.GetConfig()
	baseConfig := (*config).GetBase()
	if !baseConfig.Enable{
		return module.ErrModuleDisabled
	}
	if err := c.CheckStart();err!=nil{
		return err
	}

	//延迟连接
    newRdb := NewRedisClient(baseConfig)
	//替换
	c.setRDB(newRdb)
	c.DoStart()
    return nil
}

func (c *Client) Close(_ context.Context) error{
    c.mu.Lock()
    defer c.mu.Unlock()
	if err := c.CheckClose();err!=nil{
		return err
	}
    if c.rdb!=nil{
        _ = c.rdb.Close()
    }
	c.DoInit()
    return nil
}

func (c *Client) setRDB(newRdb redis.UniversalClient) {
	old := c.rdb
	c.rdb = newRdb

	//主动关闭
	if old!=nil{
		go func() {
			time.Sleep(2 * time.Second)
			_ = old.Close()
		}()
	}
}

//SetRDB 替换client
// 运行中才能替换
func (c *Client) SetRDB(newRdb redis.UniversalClient) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.IsRunning(){
		return module.ErrNotRunning
	}
	c.setRDB(newRdb)
	return nil
}

func (c *Client) WrapCtx(parent context.Context) (context.Context, context.CancelFunc){
	return context.WithTimeout(parent, 2*time.Second)
}

//View 访问client函数
// 1.外部禁止直接持有client
// 2.未提供的api接口,通过view访问
func (c *Client) View(fn func(rdb redis.UniversalClient) error) error{
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.IsRunning(){
		return module.ErrModuleDisabled
	}
	if c.rdb==nil{
		return ErrClientNil
	}
	return fn(c.rdb)
}

//-----------------------------------------------
//    常用接口
//-----------------------------------------------

func (c *Client) Ping(ctx context.Context) error{
	return c.View(func(rdb redis.UniversalClient) error {
		return rdb.Ping(ctx).Err()
	})
}

func (c *Client) Set(ctx context.Context, key string, val any, expire time.Duration) error {
	return c.View(func(rdb redis.UniversalClient) error {
		return rdb.Set(ctx, key, val, expire).Err()
	})
}

//SetNX 不存在才设置
func (c *Client)  SetNX(ctx context.Context, key string, val any, expire time.Duration) (rtVal bool, err error){
	err = c.View(func(rdb redis.UniversalClient) error{
		rtVal, err = rdb.SetNX(ctx, key, val, expire).Result()
		return err
	})
	return rtVal, err
}

func (c *Client) Get(ctx context.Context, key string) (rtVal string, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.Get(ctx, key).Result()
		return err
	})
	return rtVal, err
}

//GetInto dest 要求是指针
func (c *Client) GetInto(ctx context.Context, key string, dest any) error{
	return c.View(func(rdb redis.UniversalClient) error {
		return GetInto(rdb, ctx, key, dest)
	})
}


func (c *Client) Del(ctx context.Context, keys ...string) (rtVal int64, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.Del(ctx, keys...).Result()
		return err
	})
	return rtVal, err
}

//Exists 判断key是否存在
func (c *Client) Exists(ctx context.Context, keys ...string) (rtVal bool, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		var n int64
		n, err = rdb.Exists(ctx, keys...).Result()
		rtVal = n > 0
		return err
	})
	return rtVal, err
}

//Expire 设置key过期
func (c *Client) Expire(ctx context.Context, key string, expire time.Duration) (rtVal bool, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.Expire(ctx, key, expire).Result()
		return err
	})
	return rtVal, err
}

func (c *Client) TTL(ctx context.Context, key string) (rtVal time.Duration, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.TTL(ctx, key).Result()
		return err
	})
	return rtVal, err
}

//Incr 计数自增
func (c *Client) Incr(ctx context.Context, key string) (rtVal int64, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.Incr(ctx, key).Result()
		return err
	})
	return rtVal, err
}

//Decr 自减
func (c *Client) Decr(ctx context.Context, key string) (rtVal int64, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.Decr(ctx, key).Result()
		return err
	})
	return rtVal, err
}

func (c *Client) HSet(ctx context.Context, key string, fieldValues ...any) (rtVal int64, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.HSet(ctx, key, fieldValues...).Result()
		return err
	})
	return rtVal, err
}

func (c *Client) HGet(ctx context.Context, key string, field  string) (rtVal string, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.HGet(ctx, key, field).Result()
		return err
	})
	return rtVal, err
}

func (c *Client) HGetAll(ctx context.Context, key string) (rtVal map[string]string, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.HGetAll(ctx, key).Result()
		return err
	})
	return rtVal, err
}

//LPush list头部插入
func (c *Client) LPush(ctx context.Context, key string, values ...any) (rtVal int64, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.LPush(ctx, key, values...).Result()
		return err
	})
	return rtVal, err
}

//RPop list尾部弹出
func (c *Client) RPop(ctx context.Context, key string) (rtVal string, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.RPop(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	})
	return rtVal, err
}

//MGet 批量获取
func (c *Client) MGet(ctx context.Context, keys ...string) (rtVal []interface{}, err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		rtVal, err = rdb.MGet(ctx, keys...).Result()
		return err
	})
	return rtVal, err
}

//MSet 批量设置 [key, value, key, value,...]
func (c *Client) MSet(ctx context.Context, keyValues ...interface{}) (err error){
	err = c.View(func(rdb redis.UniversalClient) error {
		return rdb.MSet(ctx, keyValues).Err()
	})
	return err
}
