package main

import (
	"fmt"
	"time"

	"github.com/Prem-099/lru-cache/lru"
)

func main() {
	cache := lru.New[string, int](100)
	cache.Put("a", 1, 0*time.Second)
	cache.Put("b", 2, 0*time.Second)
	val, ok := cache.Get("a")
	fmt.Println("Value of key a : ", val, ok)
	cache.Put("c", 3, 2*time.Second)
	time.Sleep(6 * time.Second)
	_, ok = cache.Get("b")
	if ok {
		fmt.Println("b exists ", ok)
	} else {
		fmt.Println("Key expired")
	}
	_, ok = cache.Get("a")
	if ok {
		fmt.Println("a exists ", ok)
	} else {
		fmt.Println("Key expired")
	}
	stats := cache.Stats()
	fmt.Println(stats)
}
