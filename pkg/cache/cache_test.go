package cache

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}

var _ = Describe("Cache", func() {
	It("should add entries", func() {
		cache := NewCache(DefaultCacheSize, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		Expect(cache.size).To(Equal(2))
		Expect(cache.Get("key1")).To(Equal("value1"))
		Expect(cache.Get("key2")).To(Equal("value2"))
	})
	It("should delete entries", func() {
		cache := NewCache(DefaultCacheSize, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Delete("key1")

		Expect(cache.size).To(Equal(1))
		Expect(cache.Get("key1")).To(BeNil())
		Expect(cache.Get("key2")).To(Equal("value2"))
	})
	It("should add a new entry and evict the oldest (not used)", func() {
		cache := NewCache(3, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")
		cache.Set("key4", "value4")

		Expect(cache.size).To(Equal(3))
		Expect(cache.Get("key1")).To(BeNil())
		Expect(cache.Get("key2")).To(Equal("value2"))
		Expect(cache.Get("key3")).To(Equal("value3"))
		Expect(cache.Get("key4")).To(Equal("value4"))
	})
	It("should add a new entry and evict the oldest (least recently used)", func() {
		cache := NewCache(3, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")
		cache.Get("key2")
		cache.Get("key1")
		cache.Get("key3")
		cache.Set("key4", "value4")

		Expect(cache.size).To(Equal(3))
		Expect(cache.Get("key2")).To(BeNil())
		Expect(cache.Get("key1")).To(Equal("value1"))
		Expect(cache.Get("key3")).To(Equal("value3"))
		Expect(cache.Get("key4")).To(Equal("value4"))
	})
	It("should expire entries that have passed the ttl", func() {
		cache := NewCache(DefaultCacheSize, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.items["key1"].created = time.Now().Add(-time.Duration(DefaultCacheItemTTL))
		cache.Set("key2", "value2")
		cache.items["key2"].created = time.Now().Add(-time.Duration(DefaultCacheItemTTL / 2)) // not expired

		Expect(cache.size).To(Equal(2))
		Expect(cache.Get("key1")).To(BeNil())
		Expect(cache.Get("key2")).To(Equal("value2"))
		Expect(cache.size).To(Equal(1))
	})
	It("should clear the cache", func() {
		cache := NewCache(DefaultCacheSize, DefaultCacheItemTTL)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Clear()

		Expect(cache.size).To(Equal(0))
		Expect(cache.Get("key1")).To(BeNil())
		Expect(cache.Get("key2")).To(BeNil())
	})
})
