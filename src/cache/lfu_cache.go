package cache

import (
	"container/heap"
	"sync"
)

// Each key maintains a frequency count, and the cache evicts the least frequently used key when full
type LFUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex // cannot use RWMutex here because we need to update the frequency
	data     map[K]*Item[K, V]
	heap     *CacheMinHeap[K, V]
}

func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return &LFUCache[K, V]{
		data:     make(map[K]*Item[K, V]),
		heap:     &CacheMinHeap[K, V]{},
		capacity: capacity,
	}
}

func (c *LFUCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *LFUCache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.data[key]
	return exists
}

func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.data[key]; !ok {
		var missing V
		return missing, false

	} else {
		item.priority++
		heap.Fix(c.heap, item.index)
		return item.value, true
	}
}

func (c *LFUCache[K, V]) GetMany(keys []K) (map[K]V, []K) {
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

func (c *LFUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &Item[K, V]{key: key, value: value, priority: 1}
	c.data[key] = item
	heap.Push(c.heap, item)

	if len(c.data) > c.capacity {
		evicted := heap.Pop(c.heap).(*Item[K, V])
		delete(c.data, evicted.key)
	}
}

func (c *LFUCache[K, V]) SetMany(items map[K]V) {
	for key, value := range items {
		c.Set(key, value)
	}
}

func (c *LFUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.data[key]; ok {
		heap.Remove(c.heap, item.index)
		delete(c.data, key)
	}
}

func (c *LFUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[K]*Item[K, V])
	c.heap = &CacheMinHeap[K, V]{}
	heap.Init(c.heap)
}

func (c *LFUCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func (c *LFUCache[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]V, 0, len(c.data))
	for _, item := range c.data {
		values = append(values, item.value)
	}
	return values
}

var _ Cache[string, string] = (*LFUCache[string, string])(nil)
