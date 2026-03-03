package cache

import (
	"sync"
	"wbtech/internal/domain/order"
)

type OrderCache struct {
	mu       sync.RWMutex
	ordermap map[string]*order.Order
	maxsize  int
	queue    *Queue
}

func NewOrderCache(maxSize int) *OrderCache {
	return &OrderCache{
		ordermap: make(map[string]*order.Order),
		queue:    NewQueue(maxSize),
		maxsize:  maxSize,
	}
}

func (c *OrderCache) Set(key string, val *order.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.ordermap[key]; ok {
		c.ordermap[key] = val
		return
	}
	if len(c.ordermap) >= c.maxsize {
		keytodelete := c.queue.Pop().(string)
		delete(c.ordermap, keytodelete)
	}
	c.ordermap[key] = val
	c.queue.Push(val.ID)

}

func (c *OrderCache) Get(key string) *order.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.ordermap[key]; ok {
		return val
	}
	return nil
}
