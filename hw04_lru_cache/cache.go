package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*cacheItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Creates a new cache item.
func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

// Creates a new cache.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*cacheItem, capacity),
	}
}

// Sets a new value by key.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	i, ok := c.items[key]
	if ok {
		i.value = value
		c.queue.MoveToFront(c.queue.Get(i))
		return true
	}

	newItem := newCacheItem(key, value)

	c.queue.PushFront(newItem)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		c.queue.Remove(last)

		ci, ok := last.Value.(*cacheItem)
		if !ok {
			panic("failed to cast to *cacheItem")
		}

		delete(c.items, ci.key)
	}

	return false
}

// Returns the set value by key.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(c.queue.Get(val))

	return val.value, true
}

// Clears the cache.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*cacheItem)
	c.queue = NewList()
}
