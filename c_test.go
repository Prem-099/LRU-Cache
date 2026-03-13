package main

/*
import (
	"testing"
	"github.com/Prem-099/lru-cache/lru"
)
/*
func BenchmarkTest( b *testing.B){
	cache := lru.New[string,int](1000)
	cache.Put("a",1)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		cache.Get("a")
	}
}

func BenchmarkTest( b *testing.B){
	cache := lru.NewSharded[string,int](1000,16)
	cache.Put("a",1)
	b.RunParallel(func(p *testing.PB) {
		for p.Next(){
			cache.Get("a")
		}
	})
}

func BenchmarkTest(b *testing.B) {
	cache := lru.New[int,int](1000)
	b.RunParallel(func(p *testing.PB) {
		i:=0
		for p.Next(){
			cache.Put(i,i)
			cache.Get(i)
			i++
		}
	})
}

func BenchmarkWriteHeavy(b *testing.B) {
	cache := lru.New[int,int](1000)
	b.RunParallel(func(p *testing.PB) {
		i:=0
		for p.Next(){
			cache.Put(i,i)
			i++
		}
	})
}

func BenchmarkHeavyWriteShard(b *testing.B) {
	cache := lru.NewSharded[int,int](1000,256)
	b.RunParallel(func(p *testing.PB) {
		i:=0;
		for p.Next(){
			cache.Put(i,i)
			i++
		}
	})
}
/*
func BenchmarkParallelGetShard(b *testing.B) {
	cache := lru.NewSharded[string,int](1000,256)
	cache.Put("key",1)
	b.RunParallel(func(p *testing.PB) {
		for p.Next(){
			cache.Get("key")
		}
	})
}
*/