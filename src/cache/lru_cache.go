package cache

import (
	"container/heap"
	"sync"
	"time"
)

type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// Each key maintains a timestamp, and the cache evicts the least recently used key when full
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	data     map[K]*Item[K, V]
	heap     *CacheMinHeap[K, V]
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		data:     make(map[K]*Item[K, V]),
		heap:     &CacheMinHeap[K, V]{},
		capacity: capacity,
	}
}

func (c *LRUCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *LRUCache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.data[key]
	return exists
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.data[key]; !ok {
		var missing V
		return missing, false

	} else {
		item.priority = int(time.Now().UTC().Unix())
		heap.Fix(c.heap, item.index)
		return item.value, true
	}
}

func (c *LRUCache[K, V]) GetMany(keys []K) (map[K]V, []K) {
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

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &Item[K, V]{key: key, value: value, priority: int(time.Now().UTC().Unix())}
	c.data[key] = item
	heap.Push(c.heap, item)

	if len(c.data) > c.capacity {
		evicted := heap.Pop(c.heap).(*Item[K, V])
		delete(c.data, evicted.key)
	}
}

func (c *LRUCache[K, V]) SetMany(items map[K]V) {
	for key, value := range items {
		c.Set(key, value)
	}
}

func (c *LRUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.data[key]; ok {
		heap.Remove(c.heap, item.index)
		delete(c.data, key)
	}
}

func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[K]*Item[K, V])
	c.heap = &CacheMinHeap[K, V]{}
	heap.Init(c.heap)
}

func (c *LRUCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func (c *LRUCache[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]V, 0, len(c.data))
	for _, item := range c.data {
		values = append(values, item.value)
	}
	return values
}

var _ Cache[string, string] = (*LRUCache[string, string])(nil)
