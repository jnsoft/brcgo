package cache

import "sync"

// Cache is a simple in-memory cache with a mutex for thread-safe access.
type Cache struct {
	mu    sync.Mutex
	store map[string]string
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.store[key]
	return value, exists
}

// Set stores a value in the cache with the given key.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

// Delete removes a value from the cache by key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Clear removes all values from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]string)
}

// Size returns the number of items in the cache.
func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.store)
}

// Keys returns a slice of all keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]string, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice of all values in the cache.
func (c *Cache) Values() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	values := make([]string, 0, len(c.store))
	for _, value := range c.store {
		values = append(values, value)
	}
	return values
}

// Contains checks if the cache contains a specific key.
func (c *Cache) Contains(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.store[key]
	return exists
}

// Merge merges another cache into the current cache.
func (c *Cache) Merge(other *Cache) {
	c.mu.Lock()
	defer c.mu.Unlock()
	other.mu.Lock()
	defer other.mu.Unlock()
	for key, value := range other.store {
		c.store[key] = value
	}
}

// Clone creates a shallow copy of the cache.
func (c *Cache) Clone() *Cache {
	c.mu.Lock()
	defer c.mu.Unlock()
	clone := NewCache()
	for key, value := range c.store {
		clone.store[key] = value
	}
	return clone
}
