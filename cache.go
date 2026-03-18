package lru

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache[K comparable, V any] struct {
	mux         sync.Mutex
	items       map[K]*Node[K, V]
	list        *List[K, V]
	_           [56]byte
	moveCounter uint64
	pool        sync.Pool
	capacity    int
	metrics     Metrics
	stopJanitor chan struct{}
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <=0 {
		capacity = 100
	}
	c :=  &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*Node[K, V], capacity),
		list:     &List[K, V]{},
		pool: sync.Pool{
			New: func() any {
				return new(Node[K, V])
			},
		},
	}
	return c
}

var now int64

func init() {
	atomic.StoreInt64(&now, time.Now().UnixNano())
	go func() {
		for {
			atomic.StoreInt64(&now, time.Now().UnixNano())
			time.Sleep(time.Millisecond)
		}
	}()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mux.Lock()
	node, ok := c.items[key]
	if !ok {
		c.mux.Unlock()
		c.metrics.Misses.Add(1)
		var zero V
		return zero, false
	}
	if node.exp > 0 && now > node.exp {
		c.list.RemoveNode(node)
		delete(c.items, key)
		node.prev = nil
		node.next = nil
		c.pool.Put(node)
		c.mux.Unlock()
		c.metrics.Expirations.Add(1)
		var zero V
		return zero, false
	}
	c.moveCounter++
	if c.moveCounter&16==0{
		c.list.MoveToFront(node)
	}
	val := node.value
	c.mux.Unlock()
	c.metrics.Hits.Add(1)
	return val, true
}

func (c *Cache[K, V]) Put(key K, Value V, ttl time.Duration) {
	c.metrics.Puts.Add(1)
	exp := int64(0)
	c.mux.Lock()
	if ttl > 0 {
		exp = now + int64(ttl)
	}
	if node, ok := c.items[key]; ok {
		node.value = Value
		node.exp = exp
		c.list.MoveToFront(node)
		c.mux.Unlock()
		return
	}
	newnode := c.pool.Get().(*Node[K, V])
	newnode.key = key
	newnode.value = Value
	newnode.exp = exp
	c.list.AddToFront(newnode)
	c.items[key] = newnode
	if len(c.items) > c.capacity {
		lru := c.list.RemoveTail()
		if lru != nil {
			delete(c.items, lru.key)
			lru.prev = nil
			lru.next = nil
			c.pool.Put(lru)
			c.mux.Unlock()
			c.metrics.Evictions.Add(1)
			return
		}
	}
	c.mux.Unlock()
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.mux.Lock()
	node, ok := c.items[key]
	if !ok {
		c.mux.Unlock()
		return false
	}
	c.list.RemoveNode(node)
	delete(c.items, node.key)
	node.next = nil
	node.prev = nil
	c.pool.Put(node)
	c.mux.Unlock()
	return true
}

func (c *Cache[K, V]) Len() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.items)
}

func (c *Cache[K, V]) Clear() {
	c.mux.Lock()
	c.items = make(map[K]*Node[K, V])
	c.list = &List[K, V]{}
	c.mux.Unlock()
}
