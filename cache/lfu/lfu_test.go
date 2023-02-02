package lfu_test

import (
	"testing"

	"github.com/geekr-dev/go-cache/cache/lfu"
	"github.com/matryer/is"
)

func TestSetGet(t *testing.T) {
	is := is.New(t)

	cache := lfu.New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	is.Equal(v, 1)

	cache.Del("k1")
	is.Equal(0, cache.Len())
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := lfu.New(32, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	// cache.Get("k1")
	// cache.Get("k2")
	cache.Set("k3", 3)
	cache.Set("k4", 4)

	expected := []string{"k1", "k3"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}
