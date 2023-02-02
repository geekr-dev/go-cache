package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type GeekCache struct {
	cache  *safeCache // 缓存存储器
	getter Getter     // 缓存未命中时获取源数据的回调函数
}

func NewGeekCache(getter Getter, cache Cache) *GeekCache {
	return &GeekCache{
		cache:  newSafeCache(cache),
		getter: getter,
	}
}

func (g *GeekCache) Get(key string) interface{} {
	val := g.cache.get(key)
	if val != nil {
		return val
	}
	if g.getter != nil {
		val = g.getter.Get(key)
		if val == nil {
			return nil
		}
		g.cache.set(key, val)
		return val
	}
	return nil
}

func (g *GeekCache) Stat() *Stat {
	return g.cache.stat()
}
