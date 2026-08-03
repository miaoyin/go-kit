package nats

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

//---------------------------------------
//    订阅管理
//  1.支持重新订阅
//---------------------------------------

type MsgHandler nats.MsgHandler

type SubscriptionItem struct {
	Subject    string
	QueueGroup string
	Handler    MsgHandler
	sub        *nats.Subscription
}

type SubscriberManager struct {
	mu          sync.RWMutex
	subscribers map[string]*SubscriptionItem
}

func NewSubscriberManager() *SubscriberManager {
	return &SubscriberManager{
		subscribers: make(map[string]*SubscriptionItem),
	}
}

func (s *SubscriberManager) key(subject, queue string) string {
	return subject + "|" + queue
}

//Register 注册
//  在nats启动前先注册, 避免动态注册
func (s *SubscriberManager) Register(subject, queueGroup string, handler MsgHandler) error {
	if subject == "" || handler == nil {
		return errors.New("invalid subject or handler")
	}
	key := s.key(subject, queueGroup)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.subscribers[key]; exist {
		return errors.New("subscription already registered: " + key)
	}

	s.subscribers[key] = &SubscriptionItem{
		Subject:    subject,
		QueueGroup: queueGroup,
		Handler:    handler,
	}
	return nil
}

func (s *SubscriberManager) UnRegister(subject, queueGroup string) error {
	key := s.key(subject, queueGroup)
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.subscribers[key]
	if !ok {
		return errors.New("subscription not found")
	}
	if item.sub != nil {
		_ = item.sub.Unsubscribe()
	}
	delete(s.subscribers, key)
	return nil
}

func (s *SubscriberManager) ResubscribeAll(conn *nats.Conn) error {
	if conn == nil || conn.IsClosed() {
		return errors.New("invalid nats conn")
	}

	s.mu.RLock()
	list := make([]*SubscriptionItem, 0, len(s.subscribers))
	for _, v := range s.subscribers {
		list = append(list, v)
	}
	s.mu.RUnlock()

	var firstErr error
	for _, item := range list {
		var sub *nats.Subscription
		var err error

		if item.QueueGroup == "" {
			sub, err = conn.Subscribe(item.Subject, nats.MsgHandler(item.Handler))
		} else {
			sub, err = conn.QueueSubscribe(item.Subject, item.QueueGroup, nats.MsgHandler(item.Handler))
		}
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		// 保存新subscription对象
		s.mu.Lock()
		item.sub = sub
		s.mu.Unlock()
	}
	return firstErr
}

func (s *SubscriberManager) CloseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.subscribers {
		if item.sub != nil {
			_ = item.sub.Unsubscribe()
		}
	}
	s.subscribers = make(map[string]*SubscriptionItem)
}