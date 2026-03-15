package lru

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache[K comparable, V any] struct {
	capacity    int
	items       map[K]*Node[K, V]
	list        *List[K, V]
	mux         sync.Mutex
	pool        sync.Pool
	moveCounter uint64
	metrics     Metrics
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*Node[K, V], capacity),
		list:     &List[K, V]{},
		pool: sync.Pool{
			New: func() any {
				return new(Node[K, V])
			},
		},
	}
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
	if node.exp > 0 && atomic.LoadInt64(&now) > node.exp {
		c.metrics.Expirations.Add(1)
		c.list.RemoveNode(node)
		delete(c.items, key)
		c.mux.Unlock()
		var zero V
		return zero, false
	}
	c.moveCounter++
	if c.moveCounter&7 == 0 {
		c.list.MoveToFront(node)
	}
	c.mux.Unlock()
	c.metrics.Hits.Add(1)
	return node.value, true
}

func (c *Cache[K, V]) Put(key K, Value V, ttl time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.metrics.Puts.Add(1)
	exp := int64(0)
	if ttl <= 0 {
		ttl = 5 * time.Second
	}
	exp = atomic.LoadInt64(&now) + int64(ttl)
	if node, ok := c.items[key]; ok {
		node.value = Value
		node.exp = exp
		c.list.MoveToFront(node)
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
			c.metrics.Evictions.Add(1)
			lru.prev = nil
			lru.next = nil
			c.pool.Put(lru)
		}
	}
}
