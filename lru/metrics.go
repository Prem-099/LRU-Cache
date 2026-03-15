package lru

import "sync/atomic"

type Metrics struct {
	Hits        atomic.Uint64
	Misses      atomic.Uint64
	Evictions   atomic.Uint64
	Expirations atomic.Uint64
	Puts        atomic.Uint64
}

type MetricsSnapshot struct {
	Hits        uint64
	Misses      uint64
	Evictions   uint64
	Expirations uint64
	Puts        uint64
}

func (c *Cache[K, V]) Stats() MetricsSnapshot {
	return MetricsSnapshot{
		Hits:        c.metrics.Hits.Load(),
		Misses:      c.metrics.Misses.Load(),
		Evictions:   c.metrics.Evictions.Load(),
		Expirations: c.metrics.Expirations.Load(),
		Puts:        c.metrics.Puts.Load(),
	}
}

func (s *ShardedCache[K, V]) Stats() MetricsSnapshot {
	var total MetricsSnapshot
	for _, shard := range s.shards {
		stats := shard.Stats()
		total.Hits += stats.Hits
		total.Misses += stats.Misses
		total.Evictions += stats.Evictions
		total.Expirations += stats.Expirations
		total.Puts += stats.Puts
	}
	return total
}

func (m MetricsSnapshot) HitRate() float64 {
	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Hits) / float64(total)
}

func (m MetricsSnapshot) MissRate() float64 {
	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Misses) / float64(total)
}
