package cache

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

type Cache interface {
	Set(key string, val interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

type Value interface {
	Len() int
}

type safeCache struct {
	m          sync.RWMutex
	cache      Cache
	nhit, nget int
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{cache: cache}
}

func (sc *safeCache) set(key string, val interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, val)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}

	val := sc.cache.Get(key)
	if val != nil {
		log.Printf("cache hit: %s\n", key)
		sc.nhit++
	}
	return val
}

func (sc *safeCache) stat() *Stat {
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{sc.nhit, sc.nget}
}

type Stat struct {
	NHit, NGet int
}

func CalcLen(value interface{}) int {
	switch v := value.(type) {
	case string:
		if runtime.GOARCH == "amd64" {
			return 16 + len(v)
		} else {
			return 8 + len(v)
		}
	case bool, int8, uint8:
		return 1
	case int16, uint16:
		return 2
	case int32, uint32, float32:
		return 4
	case int64, uint64, float64:
		return 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			return 8
		} else {
			return 4
		}
	case complex64:
		return 8
	case complex128:
		return 16
	case Value:
		return v.Len()
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}
