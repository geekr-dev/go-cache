package lru

import (
	"container/list"

	"github.com/geekr-dev/go-cache/cache"
)

// LRU 缓存实现
// 数据结构和 FIFO 一样，只是在 Get 方法中，将访问的元素移动到链表尾部，这样就可以保证最近访问的元素总是在链表尾部不被移除
type lru struct {
	maxBytes  int
	onEvicted func(key string, value interface{})
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
}

// Del implements cache.Cache
func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.ll.Remove(e)
	kv := e.Value.(*entry)
	l.usedBytes -= kv.Len()
	delete(l.cache, kv.key)
	if l.onEvicted != nil {
		l.onEvicted(kv.key, kv.value)
	}
}

// DelOldest implements cache.Cache
func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

// Get implements cache.Cache
func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e) // 每次访问都将元素移动到链表尾部，这是和FIFO算法唯一的实现区别
		kv := e.Value.(*entry)
		return kv.value
	}
	return nil
}

// Len implements cache.Cache
func (l *lru) Len() int {
	return l.ll.Len()
}

// Set implements cache.Cache
func (l *lru) Set(key string, val interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		kv := e.Value.(*entry)
		l.usedBytes += cache.CalcLen(val) - cache.CalcLen(kv.value)
		kv.value = val
	} else {
		kv := entry{key, val}
		ele := l.ll.PushBack(&kv)
		l.cache[key] = ele
		l.usedBytes += kv.Len()
		if l.maxBytes > 0 && l.maxBytes < l.usedBytes {
			l.DelOldest()
		}
	}
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
