package hw04lrucache

type Key string
type itemsMap map[Key]*ListItem

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l lruCache) Set(key Key, value interface{}) bool {
	if item, ok := l.items[key]; ok {
		item.Value = value
		l.queue.MoveToFront(item)
		return true
	}
	if l.queue.Len() == l.capacity {
		l.queue.Remove(l.queue.Back())
	}
	l.items[key] = l.queue.PushFront(value)
	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return l.queue.Front().Value, true
	}
	return nil, false
}

func (l lruCache) Clear() {
	l.queue = NewList()
	l.items = make(itemsMap, l.capacity)
}

func NewCache(capacity int) Cache {
	if capacity < 1 {
		panic("capacity > 0 required")
	}

	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(itemsMap, capacity),
	}
}
