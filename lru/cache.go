package lru

import "sync"

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*Node[K, V]
	list     *List[K, V]
	mux      sync.Mutex
	pool     sync.Pool
}

func New[K comparable,V any] (capacity int) *Cache[K,V] {
	return &Cache[K,V]{
		capacity: capacity,
		items: make(map[K]*Node[K, V],capacity),
		list: &List[K, V]{},
		pool: sync.Pool{
			New: func() any {
				return new(Node[K,V])
			},
		},

	}
}

func (c *Cache[K, V]) Get(Key K) (V, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if node,ok := c.items[Key]; ok{
		c.list.MoveToFront(node)
		return node.value,true
	}
	var zero V
	return zero,false
}

func (c *Cache[K, V]) Put(Key K, Value V) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if node,ok := c.items[Key];ok{
		node.value = Value
		c.list.MoveToFront(node)
		return
	}
	/*newnode := &Node[K,V]{
		key: Key,
		value: Value,
	}*/
	newnode:=c.pool.Get().(*Node[K,V])
	newnode.key = Key
	newnode.value = Value
	c.list.AddToFront(newnode)
	c.items[Key] = newnode
	if len(c.items) > c.capacity{
		lru := c.list.RemoveTail()
		if lru!=nil{
			delete(c.items,lru.key)
			lru.prev = nil
			lru.next = nil
			c.pool.Put(lru)
		}
	}
}