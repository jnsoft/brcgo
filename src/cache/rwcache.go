package cache

import "sync"

// Cache replacement policies:
// 1. LRU (Least Recently Used)
// 2. LFU (Least Frequently Used)
// 3. FIFO (First In First Out)
// 4. Random Replacement
// 5. Clock
// 6. Cache with expiration
// 7. Cache with size limit
// 8. Hybrid policies

type Rwcache struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewCache creates a new Cache instance.
func NewRwcache() *Rwcache {
	return &Rwcache{
		store: make(map[string]string),
	}
}

func (c *Rwcache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.store[key]
	return value, exists
}

func (c *Rwcache) GetMany(keys []string) map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, exists := c.store[key]; exists {
			values[key] = value
		}
	}
	return values
}

func (c *Rwcache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Rwcache) SetMany(items map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range items {
		c.store[key] = value
	}
}

// Delete removes a value from the cache by key.
func (c *Rwcache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Clear removes all values from the cache.
func (c *Rwcache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]string)
}

// Size returns the number of items in the cache.
func (c *Rwcache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// Keys returns a slice of all keys in the cache.
func (c *Rwcache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice of all values in the cache.
func (c *Rwcache) Values() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]string, 0, len(c.store))
	for _, value := range c.store {
		values = append(values, value)
	}
	return values
}

// Contains checks if the cache contains a specific key.
func (c *Rwcache) Contains(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.store[key]
	return exists
}

// Merge merges another cache into the current cache.
func (c *Rwcache) Merge(other *Rwcache) {
	c.mu.Lock()
	defer c.mu.Unlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for key, value := range other.store {
		c.store[key] = value
	}
}

// Clone creates a shallow copy of the cache.
func (c *Rwcache) Clone() *Rwcache {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := NewRwcache()
	for key, value := range c.store {
		clone.store[key] = value
	}
	return clone
}
