package lfu

import (
	"container/heap"

	"github.com/geekr-dev/go-cache/cache"
)

type lfu struct {
	maxBytes  int
	onEvicted func(key string, value interface{})
	usedBytes int
	queue     *queue
	cache     map[string]*entry
}

// Del implements cache.Cache
func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}

	e := x.(*entry)
	delete(l.cache, e.key)
	l.usedBytes -= e.Len()
	if l.onEvicted != nil {
		l.onEvicted(e.key, e.value)
	}
}

// DelOldest implements cache.Cache
func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

// Get implements cache.Cache
func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}
	return nil
}

// Len implements cache.Cache
func (l *lfu) Len() int {
	return l.queue.Len()
}

// Set implements cache.Cache
func (l *lfu) Set(key string, val interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes += cache.CalcLen(val) - cache.CalcLen(e.value)
		l.queue.update(e, val, e.weight+1)
	} else {
		e := &entry{
			key:   key,
			value: val,
		}
		heap.Push(l.queue, e)
		l.cache[key] = e
		l.usedBytes += e.Len()
		if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
			l.removeElement(heap.Pop(l.queue))
		}
	}
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &queue{},
		cache:     make(map[string]*entry),
	}
}

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

// 使用小顶堆实现LFU队列
type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	e := x.(*entry)
	e.index = n
	*q = append(*q, e)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	e := old[n-1]
	e.index = -1
	*q = old[0 : n-1]
	return e
}

func (q *queue) update(e *entry, value interface{}, weight int) {
	e.value = value
	e.weight = weight
	heap.Fix(q, e.index)
}
