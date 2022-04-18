package cache

import (
	"sync"
	"time"
)

const (
	DefaultCacheSize    = 1024
	DefaultCacheItemTTL = 60000000000 // 1 minute in nanoseconds
)

type entry struct {
	key     string
	value   interface{}
	created time.Time
	next    *entry
	prev    *entry
}

// Cache is a basic thread-safe in-memory LRU cache.
// It is thread-safe and is used as the default items if providers
type Cache struct {
	// The maximum size of the cache.
	MaxCacheSize int
	// The time-to-live for cache entries in nanoseconds.
	CacheEntryTTL time.Duration

	size  int
	items map[string]*entry
	lock  sync.Mutex
	head  *entry
	tail  *entry
}

// New creates a new cache with default settings.
func New() *Cache {
	return NewCache(DefaultCacheSize, DefaultCacheItemTTL)
}

// NewCache creates a new cache with the specified max size and entry TTL.
func NewCache(maxCacheSize int, cacheEntryTTL time.Duration) *Cache {
	return &Cache{
		MaxCacheSize:  maxCacheSize,
		CacheEntryTTL: cacheEntryTTL,

		size:  0,
		items: make(map[string]*entry, maxCacheSize),
		lock:  sync.Mutex{},
		head:  nil,
		tail:  nil,
	}
}

func (c *Cache) Get(key string) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.items[key]
	if !ok {
		return nil
	}

	if time.Now().Sub(entry.created) > c.CacheEntryTTL {
		c.remove(entry)
		return nil
	}

	c.promote(entry)
	return entry.value
}

func (c *Cache) Set(key string, value interface{}) {
	if value == nil {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	if entry, ok := c.items[key]; ok {
		c.remove(entry)
	} else if c.size >= c.MaxCacheSize {
		c.evict()
	}

	entry := &entry{
		key:     key,
		value:   value,
		created: time.Now(),
	}
	c.add(entry)
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.items[key]
	if !ok {
		return
	}

	c.remove(entry)
}

func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = make(map[string]*entry, c.MaxCacheSize)
	c.size = 0
	c.head = nil
	c.tail = nil
}

//

func (c *Cache) add(entry *entry) {
	c.items[entry.key] = entry
	c.size++

	if c.head == nil {
		c.head = entry
		c.tail = entry
		return
	}

	entry.next = c.head
	c.head.prev = entry
	c.head = entry
}

func (c *Cache) remove(entry *entry) {
	if entry.prev != nil {
		entry.prev.next = entry.next
	} else {
		c.head = entry.next
	}

	if entry.next != nil {
		entry.next.prev = entry.prev
	} else {
		c.tail = entry.prev
	}

	entry.next = nil
	entry.prev = nil
	delete(c.items, entry.key)
	c.size--
}

func (c *Cache) promote(entry *entry) {
	if (c.size == 1) || (entry == c.head) {
		return
	}

	if entry.prev != nil {
		entry.prev.next = entry.next
	} else {
		c.head = entry.next
	}

	if entry.next != nil {
		entry.next.prev = entry.prev
	} else {
		c.tail = entry.prev
	}

	entry.next = c.head
	c.head.prev = entry
	c.head = entry
}

func (c *Cache) evict() {
	if c.tail == nil {
		return
	}
	c.remove(c.tail)
}
