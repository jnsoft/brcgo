package cache

import (
	"container/heap"
	"sync"
	"time"
)

// DLFU combines LRU and LFU.

type DLFUCache[K comparable, V any] struct {
	capacity int
	gamma    float64
	incr     float64
	decay    float64
	mu       sync.Mutex
	data     map[K]*DLFUItem[K, V]
	heap     *DLFUMinHeap[K, V]
}

type DLFUItem[K comparable, V any] struct {
	key       K
	value     V
	priority  float64
	expiresAt time.Time
	index     int // Index in the heap (updated by heap)
}

func (i *DLFUItem[K, V]) expired() bool {
	return i.expiresAt.Before(time.Now())
}

type DLFUMinHeap[K comparable, V any] []*DLFUItem[K, V]

func (pq DLFUMinHeap[K, V]) Len() int { return len(pq) }

func (pq DLFUMinHeap[K, V]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq DLFUMinHeap[K, V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *DLFUMinHeap[K, V]) Push(x any) {
	entry := x.(*DLFUItem[K, V])
	entry.index = len(*pq)
	*pq = append(*pq, entry)
}

func (pq *DLFUMinHeap[K, V]) Pop() any {
	old := *pq
	n := len(old)
	entry := old[n-1]
	entry.index = -1 // Mark as removed
	*pq = old[:n-1]
	return entry
}

func NewDLFUCache[K comparable, V any](capacity int, gamma float64) *DLFUCache[K, V] {
	if gamma < 0.0 || gamma > 1.0 {
		panic("gamma must be between 0 and 1")
	}

	cache := &DLFUCache[K, V]{
		data:     make(map[K]*DLFUItem[K, V]),
		heap:     &DLFUMinHeap[K, V]{},
		capacity: capacity,
		gamma:    gamma,
	}

	heap.Init(cache.heap)

	if gamma < 1.0 {
		cache.incr = 1.0 / (1.0 - gamma)
	}

	p := float64(capacity) * gamma
	cache.decay = (p + 1.0) / p

	return cache
}

func (c *DLFUCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *DLFUCache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.data[key]
	return exists
}

func (c *DLFUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if ok && !item.expired() {
		item.priority += c.incr
		heap.Fix(c.heap, item.index)
		c.incr *= c.decay
		return item.value, true

	} else {
		c.incr *= c.decay
		var missing V
		return missing, false
	}
}

func (c *DLFUCache[K, V]) GetMany(keys []K) (map[K]V, []K) {
	values := make(map[K]V, len(keys))
	missingKeys := make([]K, 0)
	for _, key := range keys {
		if item, ok := c.Get(key); ok {
			values[key] = item
		} else {
			missingKeys = append(missingKeys, key)
		}
	}
	return values, missingKeys
}

func (c *DLFUCache[K, V]) Set(key K, value V, expiry time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiryTime := time.Now().Add(expiry)

	// Update existing item
	if item, exists := c.data[key]; exists {
		item.value = value
		item.expiresAt = expiryTime
		item.priority = c.incr
		heap.Fix(c.heap, item.index)
		return
	}

	// Add new item
	item := &DLFUItem[K, V]{key: key, value: value, priority: c.incr, expiresAt: expiryTime}
	c.data[key] = item
	heap.Push(c.heap, item)

	// Remove expired items
	for c.heap.Len() > 0 {
		top := (*c.heap)[0]
		if top.expired() {
			heap.Pop(c.heap)
			delete(c.data, top.key)
		} else {
			break
		}
	}

	// Evict if over capacity
	if len(c.data) > c.capacity {
		evicted := heap.Pop(c.heap).(*DLFUItem[K, V])
		delete(c.data, evicted.key)
	}
}

func (c *DLFUCache[K, V]) SetMany(items map[K]V, expiry time.Duration) {
	for key, value := range items {
		c.Set(key, value, expiry)
	}
}

func (c *DLFUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.data[key]; ok {
		heap.Remove(c.heap, item.index)
		delete(c.data, key)
	}
}

func (c *DLFUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[K]*DLFUItem[K, V])
	c.heap = &DLFUMinHeap[K, V]{}
	heap.Init(c.heap)
}

func (c *DLFUCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func (c *DLFUCache[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]V, 0, len(c.data))
	for _, item := range c.data {
		values = append(values, item.value)
	}
	return values
}

func (c *DLFUCache[K, V]) removeExpiredItems() {
	for c.heap.Len() > 0 {
		top := (*c.heap)[0]
		if top.expired() {
			heap.Pop(c.heap)
			delete(c.data, top.key)
		} else {
			break
		}
	}

	// Rebuild to remove items deeper in the heap
	newHeap := &DLFUMinHeap[K, V]{}
	for _, item := range *c.heap {
		if !item.expired() {
			*newHeap = append(*newHeap, item)
		} else {
			delete(c.data, item.key)
		}
	}
	heap.Init(newHeap)
	c.heap = newHeap
}
