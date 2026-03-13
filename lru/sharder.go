package lru

type ShardedCache[K comparable, V any] struct {
	shards []*Cache[K, V]
	mask   uint64
}

func NewSharded[K comparable, V any](capacity int, shardCount int) *ShardedCache[K,V] {
	if shardCount&(shardCount-1)!=0{
		panic("ShardCount should be power of 2")
	}
	shards := make([]*Cache[K, V], shardCount)
	perShard := capacity / shardCount // capacity of each shard
	for i := 0; i < shardCount; i++ {
		shards[i] = New[K, V](perShard)
	}
	return &ShardedCache[K, V]{
		shards: shards,
		mask:   uint64(shardCount - 1), // shardCount should be power of 2 for better performance and efficiency
	}
}

func (s *ShardedCache[K, V]) getShard(key K) *Cache[K,V]{
	index := hashKey(key) & s.mask
	return s.shards[index]
}

func hashKey[K comparable](key K) uint64{
	switch v:= any(key).(type){
	case int:
		return uint64(v)
	case uint64:
		return v
	case string:
		var h uint64 = 1469598103934665603
		for i:=0;i<len(v);i++{
			h^=uint64(v[i])
			h*=1099511628211
		}
		return h
	default:
		panic("Unsupported key type")
	}
}

func (s *ShardedCache[K, V]) Get(key K) (V,bool) {
	shard := s.getShard(key)
	return shard.Get(key)
}

func (s *ShardedCache[K, V]) Put(key K,value V) {
	shard := s.getShard(key)
	shard.Put(key,value)
}