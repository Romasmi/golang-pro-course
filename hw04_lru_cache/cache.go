package hw04lrucache

type Key string
type itemsMap map[Key]*ListItem
type keyValue struct {
	key   Key
	value interface{}
}

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
	kv := &keyValue{
		key:   key,
		value: value,
	}
	if item, ok := l.items[key]; ok {
		item.Value = kv
		l.queue.MoveToFront(item)
		return true
	}
	l.items[key] = l.queue.PushFront(kv)
	if l.queue.Len() <= l.capacity {
		return false
	}
	if last := l.queue.Back(); last != nil {
		kv, _ := last.Value.(keyValue)
		delete(l.items, kv.key)
		l.queue.Remove(last)
	}
	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		kv, _ := item.Value.(*keyValue)
		return kv.value, true
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
