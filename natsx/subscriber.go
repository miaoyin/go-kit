package natsx

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

type SubDef struct {
	Subject string
	Queue   string
	Handler nats.MsgHandler
}

type SubscriberManager struct {
	mu sync.RWMutex
	//订阅信息
	defs map[string]*SubDef
	//订阅连接
	subs map[string]*nats.Subscription
}

func NewSubscriberManager() *SubscriberManager {
	sm := &SubscriberManager{
		defs: make(map[string]*SubDef),
	}
	return sm
}

func (sm *SubscriberManager) key(subject, queue string) string {
	return subject + "|" + queue
}

// RegisterSub 注册信息
// SubscribeAll 执行真实注册
func (sm *SubscriberManager) RegisterSub(subject, queueGroup string, handler nats.MsgHandler) error {
	if subject == "" || handler == nil {
		return errors.New("invalid subject or handler")
	}
	key := sm.key(subject, queueGroup)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.defs[key]; exists {
		return errors.New("subscription already registered: " + key)
	}

	def := &SubDef{
		Subject: subject,
		Queue:   queueGroup,
		Handler: handler,
	}
	sm.defs[key] = def
	return nil
}

func (sm *SubscriberManager) doSubscribe(nc *nats.Conn, def *SubDef) (*nats.Subscription, error) {
	if def.Queue == "" {
		return nc.Subscribe(def.Subject, def.Handler)
	} else {
		return nc.QueueSubscribe(def.Subject, def.Queue, def.Handler)
	}
}

func (sm *SubscriberManager) UnRegister(subject, queueGroup string) error {
	key := sm.key(subject, queueGroup)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.defs[key]; !ok {
		return errors.New("subscription not found")
	}

	if sm.subs != nil {
		sub, _ := sm.subs[key]
		if sub != nil {
			_ = sub.Unsubscribe()
		}
		delete(sm.subs, key)
	}

	delete(sm.defs, key)
	return nil
}

// RebuildAll 订阅
func (sm *SubscriberManager) RebuildAll(nc *nats.Conn) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 1. 先取消所有旧订阅
	for _, sub := range sm.subs {
		_ = sub.Unsubscribe()
	}
	sm.subs = make(map[string]*nats.Subscription)

	// 2. 使用新连接快照重建全部订阅
	for key, def := range sm.defs {
		sub, err := sm.doSubscribe(nc, def)
		if err != nil {
			return err
		}
		sm.subs[key] = sub
	}
	return nil
}

func (sm *SubscriberManager) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, sub := range sm.subs {
		_ = sub.Unsubscribe()
	}
	sm.subs = nil
}
