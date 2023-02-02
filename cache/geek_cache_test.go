package cache_test

import (
	"log"
	"sync"
	"testing"

	"github.com/geekr-dev/go-cache/cache"
	"github.com/geekr-dev/go-cache/cache/lru"
	"github.com/matryer/is"
)

func TestCacheGet(t *testing.T) {
	db := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	// 缓存未命中时，从数据库中获取
	getter := cache.GetFunc(func(key string) interface{} {
		log.Println("[From DB] find key", key)
		if v, ok := db[key]; ok {
			return v
		}
		return nil
	})
	geekCache := cache.NewGeekCache(getter, lru.New(0, nil))

	is := is.New(t)

	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			// 对于同一个key获取两次，以尽可能命中缓存
			is.Equal(geekCache.Get(k), v)
			is.Equal(geekCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	is.Equal(geekCache.Get("unknown"), nil)
	is.Equal(geekCache.Get("unknown"), nil)

	is.Equal(geekCache.Stat().NGet, 12)
	is.Equal(geekCache.Stat().NHit, 5)
}
