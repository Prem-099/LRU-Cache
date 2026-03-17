![Go Version](https://img.shields.io/badge/Go-1.18+-blue)
# ⚡ High-Performance LRU Cache in Go

A production-grade **thread-safe LRU (Least Recently Used) cache** written in Go with:

-  O(1) operations
-  TTL (expiration support)
-  Concurrency-safe design
-  Sharded cache for high throughput
-  Built-in metrics (Atomic)
-  Memory optimization using `sync.Pool`
-  Background janitor for cleanup
-  Extensive benchmark suite
  
---

## 🚀 Features

### Core
- O(1) `Get`, `Put`, `Delete`
- Doubly linked list + hashmap design
- Generic (`K comparable, V any`)

### TTL Support
- Per-key expiration
- Lazy eviction on access
- Optional background cleanup (Janitor)

### Concurrency
- Thread-safe with `sync.Mutex`
- **Sharded cache** reduces lock contention

### Metrics
- Hits / Misses
- Evictions / Expirations
- Hit rate / Miss rate

### Memory Optimization
- `sync.Pool` for node reuse
- Zero allocations per operation

---

## Architecture

- **HashMap → O(1) lookup**
- **Doubly Linked List → LRU order**
  - Head → Most Recently Used
  - Tail → Least Recently Used

---
## Installation

### Using `go get`

```bash
go get github.com/Prem-099/lru-cache@v1.0.1
```
### Import
```bash
import "github.com/Prem-099/lru-cache"
```
### Usage

### Basic

```go
cache := lru.New[string,int](100)

cache.Put("a", 1, 0)

val, ok := cache.Get("a")

//With TTL
cache.Put("key", 100, 2*time.Second)

//Sharded Cache
sharded := lru.NewSharded[string, int](1000, 16)

sharded.Put("a", 1, 0)
val, ok := sharded.Get("a")

//⚠️ shardCount must be a power of 2
```
### Metrics
```go
stats := cache.Stats()
fmt.Println(stats.HitRate())
//Janitor
cache.StartJanitor()
defer cache.StopJanitor()
```
---
## Project Structure
```bash
lru/
├── cache.go
├── list.go
├── node.go
├── sharder.go
├── metrics.go
├── janitor.go
├── c_test.go
```
---
## 📊 Benchmark Results

### Environment:

- CPU: Intel i5-13420H (12 threads)
- OS: Windows
- Command: go test -bench=. -benchmem -cpu=1,2,4,8,12


---

### 🔹 Single-thread Get Performance

| CPUs | Cache Get (ns/op) | Map Get (ns/op) |
|------|------------------|-----------------|
| 1    | 35.14            | 6.56            |
| 2    | 24.78            | 6.41            |
| 4    | 25.28            | 6.69            |
| 8    | 24.50            | 6.62            |
| 12   | 24.69            | 6.74            |

---

### 🔹 Parallel Read Performance

| CPUs | Normal Cache (ns/op) | Sharded Cache (ns/op) |
|------|---------------------|-----------------------|
| 1    | 26.36               | 26.50                 |
| 2    | 33.74               | 55.02                 |
| 4    | 59.46               | 65.63                 |
| 8    | 80.10               | 64.20                 |
| 12   | 82.73               | 69.33                 |

---

### 🔹 Mixed Workload (Read + Write)

| CPUs | Normal Cache (ns/op) | Sharded Cache (ns/op) |
|------|---------------------|-----------------------|
| 1    | 117.2               | 113.4                 |
| 2    | 229.4               | 144.3                 |
| 4    | 256.9               | 128.7                 |
| 8    | 267.6               | 122.1                 |
| 12   | 270.9               | 118.5                 |

---

### 🔹 Write Heavy

| CPUs | Normal Cache (ns/op) | Sharded Cache (ns/op) |
|------|---------------------|-----------------------|
| 1    | 95.47               | 90.07                 |
| 2    | 170.9               | 96.36                 |
| 4    | 195.3               | 76.30                 |
| 8    | 199.8               | 73.09                 |
| 12   | 195.6               | 56.07                 |

---

### 🔹 Eviction Performance

| CPUs | LRU Eviction (ns/op) | Parallel Eviction (ns/op) |
|------|---------------------|----------------------------|
| 1    | 99.72               | 87.21                      |
| 2    | 98.12               | 105.0                      |
| 4    | 97.25               | 88.16                      |
| 8    | 95.79               | 71.37                      |
| 12   | 95.27               | 56.38                      |

---

### 🔹 Zipf Distribution (Real-world Pattern)

| CPUs | Zipf Mixed (ns/op) |
|------|-------------------|
| 1    | 76.67             |
| 2    | 63.89             |
| 4    | 71.15             |
| 8    | 60.93             |
| 12   | 66.94             |

---

### 🔹 Zipf Parallel Get

| CPUs | Performance (ns/op) |
|------|---------------------|
| 1    | 26.53               |
| 2    | 40.15               |
| 4    | 80.17               |
| 8    | 107.2               |
| 12   | 109.3               |

---

### 🔹 Hit vs Miss (50/50)

| CPUs | Performance (ns/op) |
|------|---------------------|
| 1    | 32.34               |
| 2    | 24.02               |
| 4    | 23.63               |
| 8    | 23.83               |
| 12   | 23.67               |

---

### ⚔️ Performance Comparison (Approximate)
> ⚠️ Note: Ristretto and BigCache benchmarks are approximate ranges based on public benchmarks and may vary depending on workload and hardware.

| Scenario | This Cache | Ristretto | BigCache |
|----------|-----------|-----------|----------|
| Single Get | ~24–35 ns | ~10–20 ns | ~15–25 ns |
| Parallel Get (high CPU) | ~69–82 ns | ~20–40 ns | ~25–50 ns |
| Mixed Workload | ~113–270 ns | ~50–120 ns | ~60–150 ns |
| Write Heavy | ~56–199 ns | ~40–100 ns | ~50–120 ns |
| Eviction | ~95–100 ns | ~50–80 ns | ~60–100 ns |
| Zipf Workload | ~60–76 ns | ~30–60 ns | ~40–80 ns |
| Allocations | 0 alloc/op | ~0–1 alloc/op | ~0 alloc/op |

--- 

### Key Observations

-  **0 allocations per operation**
-  **Sharded cache scales significantly better**
-  Performance degrades gracefully under concurrency
-  Stable under real-world workloads (Zipf)
-  Suitable for high-throughput systems

### Comparison with Popular Go Caches
| Feature | This Project | Ristretto | BigCache |
|--------|-------------|----------|----------|
| Eviction Policy | LRU | TinyLFU | FIFO-like |
| TTL Support | ✅ | ✅ | ✅ |
| Sharding | ✅ | Internal | Internal |
| Metrics | ✅ | ✅ | Limited |
| Generics Support | ✅ | ❌ | ❌ |
| Memory Optimization | sync.Pool | Advanced (cost-based) | Byte-based |
| Allocations | 0 alloc/op | Very low | Very low |
| Concurrency | High | Very High | Very High |
| Accuracy | Lazy LRU | Probabilistic | Approximate |

---
### Use This Cache When You Want:

- true LRU behavior
- simple + predictable eviction
- full control / learning
- clean design + readability
---

### Complexity
- Operation	Time
- Get	O(1)
- Put	O(1)
- Delete	O(1)
---
### Use Cases

- API caching
- Database query caching
- Backend performance optimization
- In-memory key-value store
---
### ⚠️ Limitations
- Sharding supports limited key types (int, string, uint64)
- Not distributed
- In-memory only
---
### Future Improvements

- Lock-free design
- Distributed cache
- Persistent storage
- Prometheus integration
---
### Contributing
Contributions are welcome 🚀

### Author
PREM CHANDU PALIVELA

### License
MIT License
