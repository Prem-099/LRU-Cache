package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Prem-099/lru-cache/lru"
)

func BenchmarkTestGet(b *testing.B) {
	cache := lru.New[int, int](1000)
	cache.Put(1, 1, 2*time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(1)
	}
}

func BenchmarkTestParallelGet(b *testing.B) {
	cache := lru.NewSharded[string, int](1000, 64)
	cache.Put("a", 1, 2*time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			cache.Get("a")
		}
	})
}

func BenchmarkParallelGetShard(b *testing.B) {
	cache := lru.NewSharded[int, int](1000, 64)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i, 2*time.Second)
	}
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Get(i % 1000)
			i++
		}
	})
}

func BenchmarkTestMixed(b *testing.B) {
	cache := lru.New[int, int](1000)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			cache.Get(i)
			i++
		}
	})
}

func BenchmarkTestMixedShard(b *testing.B) {
	cache := lru.NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			cache.Get(i)
			i++
		}
	})
}

func BenchmarkWriteHeavy(b *testing.B) {
	cache := lru.New[int, int](1000)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			i++
		}
	})
}

func BenchmarkHeavyWriteShard(b *testing.B) {
	cache := lru.NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			cache.Put(i, i, 2*time.Second)
			i++
		}
	})
}

func BenchmarkTestTtl(b *testing.B) {
	cache := lru.New[int, int](1000)
	cache.Put(1, 1, 3*time.Second)
	time.Sleep(2 * time.Second)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d hitrate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(1)
	}
}

func BenchmarkCachedMixed(b *testing.B) {
	cache := lru.New[int, int](1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i, 5*time.Second)
	}
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := i % 1500
		if i%3 == 0 {
			cache.Put(key, key, 5*time.Second)
		} else {
			cache.Get(key)
		}
	}
}

func BenchmarkZipfCache(b *testing.B) {
	cache := lru.NewSharded[int, int](1000, 64)
	b.Cleanup(func() {
		stats := cache.Stats()
		b.Logf("Cache stats : Hits:%d Misses:%d Exp:%d Evic:%d Puts:%d HitRate:%f", stats.Hits, stats.Misses,
			stats.Expirations, stats.Evictions, stats.Puts, stats.HitRate())
	})
	b.RunParallel(func(p *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		zipf := rand.NewZipf(r, 1.2, 1, 1500)
		for p.Next() {
			key := int(zipf.Uint64())
			if key%4 == 0 {
				cache.Put(key, key, 2*time.Second)
			} else {
				cache.Get(key)
			}
		}
	})
}
