package cache

import "sync"

type Cache[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
	SetMany(items map[K]V)
	GetMany(keys []K) map[K]V
	Delete(key K)
	Clear()
	Size() int
	Keys() []K
	Values() []V
	Contains(key K) bool
}

var _ Cache[string, string] = (*SimpleCache[string, string])(nil)

type SimpleCache[K comparable, V any] struct {
	mu    sync.Mutex
	store map[K]V
}

func NewSimpleCache[K comparable, V any]() *SimpleCache[K, V] {
	return &SimpleCache[K, V]{
		store: make(map[K]V),
	}
}

func (c *SimpleCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.store[key]
	return value, exists
}

func (c *SimpleCache[K, V]) GetMany(keys []K) map[K]V {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make(map[K]V, len(keys))
	for _, key := range keys {
		if value, exists := c.store[key]; exists {
			values[key] = value
		}
	}
	return values
}
func (c *SimpleCache[K, V]) SetMany(items map[K]V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range items {
		c.store[key] = value
	}
}

func (c *SimpleCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
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

func (c *SimpleCache[K, V]) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.store)
}

func (c *SimpleCache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}

func (c *SimpleCache[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]V, 0, len(c.store))
	for _, value := range c.store {
		values = append(values, value)
	}
	return values
}

func (c *SimpleCache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.store[key]
	return exists
}
