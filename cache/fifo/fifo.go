package fifo

import (
	"container/list"

	"github.com/geekr-dev/go-cache/cache"
)

// FIFO 缓存实现
type fifo struct {
	// 最大容量
	maxBytes int
	// 当一个 entry 从缓存中移除时调用该 callback 函数，默认为 nil
	onEvicted func(key string, value interface{})
	// 已使用容量
	usedBytes int
	// 通过双向链表存储值
	ll *list.List
	// 通过 map 存储键值对
	cache map[string]*list.Element
}

// Del implements cache.Cache
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	kv := e.Value.(*entry)
	f.usedBytes -= kv.Len()
	delete(f.cache, kv.key)
	if f.onEvicted != nil {
		f.onEvicted(kv.key, kv.value)
	}
}

// DelOldest implements cache.Cache
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

// Get implements cache.Cache
func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		kv := e.Value.(*entry)
		return kv.value
	}
	return nil
}

// Len implements cache.Cache
func (f *fifo) Len() int {
	return f.ll.Len()
}

// Set implements cache.Cache
func (f *fifo) Set(key string, val interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		kv := e.Value.(*entry)
		f.usedBytes += cache.CalcLen(val) - cache.CalcLen(kv.value)
		kv.value = val
	} else {
		kv := entry{key, val}
		ele := f.ll.PushBack(&kv)
		f.cache[key] = ele
		f.usedBytes += kv.Len()
		if f.maxBytes > 0 && f.maxBytes < f.usedBytes {
			f.DelOldest()
		}
	}
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &fifo{
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
