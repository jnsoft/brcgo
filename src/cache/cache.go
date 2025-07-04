package cache

// Cache replacement policies:
// 1. LRU (Least Recently Used)
// 2. LFU (Least Frequently Used)
// 3. FIFO (First In First Out)
// 4. Random Replacement
// 5. Clock
// 6. Cache with expiration
// 7. Cache with size limit
// 8. Hybrid policies

import (
	"sync"
)

type Cache[K comparable, V any] interface {
	Size() int
	Contains(key K) bool
	Get(key K) (value V, ok bool)
	GetMany(keys []K) (map[K]V, []K)
	Set(key K, value V)
	SetMany(items map[K]V)
	Delete(key K)
	Clear()
	Keys() []K
	Values() []V
}

type SimpleCache[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
}

func NewSimpleCache[K comparable, V any]() *SimpleCache[K, V] {
	return &SimpleCache[K, V]{
		store: make(map[K]V),
	}
}

func (c *SimpleCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

func (c *SimpleCache[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.store[key]
	return exists
}

func (c *SimpleCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.store[key]
	return value, exists
}

func (c *SimpleCache[K, V]) GetMany(keys []K) (map[K]V, []K) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make(map[K]V, len(keys))
	missingKeys := make([]K, 0)
	for _, key := range keys {
		if value, exists := c.store[key]; exists {
			values[key] = value
		} else {
			missingKeys = append(missingKeys, key)
		}
	}
	return values, missingKeys
}

func (c *SimpleCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *SimpleCache[K, V]) SetMany(items map[K]V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range items {
		c.store[key] = value
	}
}

func (c *SimpleCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *SimpleCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[K]V)
}

func (c *SimpleCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}

func (c *SimpleCache[K, V]) Values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]V, 0, len(c.store))
	for _, value := range c.store {
		values = append(values, value)
	}
	return values
}

var _ Cache[string, string] = (*SimpleCache[string, string])(nil)
