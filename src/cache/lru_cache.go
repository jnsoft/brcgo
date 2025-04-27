package cache

import (
	"sync"
)

type Node[K comparable, V any] struct {
	key   K
	value V
	prev  *Node[K, V]
	next  *Node[K, V]
}

// Each key maintains a timestamp, and the cache evicts the least recently used key when full
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	data     map[K]*Node[K, V]
	head     *Node[K, V]
	tail     *Node[K, V]
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		data:     make(map[K]*Node[K, V]),
		capacity: capacity,
		head:     nil,
		tail:     nil,
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

	if node, ok := c.data[key]; ok {
		c.moveToHead(node)
		return node.value, true
	}
	var zero V
	return zero, false
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

	if node, exists := c.data[key]; exists {
		node.value = value
		c.moveToHead(node)
	} else {
		newNode := &Node[K, V]{key: key, value: value}
		if len(c.data) == c.capacity {
			delete(c.data, c.tail.key)
			c.removeNode(c.tail)
		}
		c.addNode(newNode)
		c.data[key] = newNode
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

	if node, ok := c.data[key]; ok {
		c.removeNode(node)
		delete(c.data, key)
	}
}

func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[K]*Node[K, V])
	c.head = nil
	c.tail = nil
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

func (c *LRUCache[K, V]) moveToHead(node *Node[K, V]) {
	c.removeNode(node)
	c.addNode(node)
}

func (c *LRUCache[K, V]) addNode(node *Node[K, V]) {
	node.next = c.head
	node.prev = nil
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache[K, V]) removeNode(node *Node[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

var _ Cache[string, string] = (*LRUCache[string, string])(nil)
