package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value    interface{}
	Next     *ListItem
	Prev     *ListItem
	position int
}

type list struct {
	hash hashMap
}

type hashMap map[int]*ListItem

func (l list) Len() int {
	return len(l.hash)
}

func (l list) Front() *ListItem {
	return l.at(0)
}

func (l list) Back() *ListItem {
	return l.at(l.Len() - 1)
}

func (l list) PushFront(v interface{}) *ListItem {
	return l.pushTo(v, 0)
}

func (l list) PushBack(v interface{}) *ListItem {
	return l.pushTo(v, l.Len())
}

func (l list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next = i.Prev
	}
	l.decIndex(i.position)
}

func (l list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{
		hash: make(hashMap),
	}
}

func (l list) at(i int) *ListItem {
	if l.Len() == 0 {
		return nil
	}
	item, ok := l.hash[i]
	if ok == false {
		return nil
	}
	item.position = i
	return item
}

func (l list) incIndex(start int) {
	for i := l.Len(); i > start; i-- {
		l.hash[i] = l.hash[i-1]
	}
}

func (l list) decIndex(start int) {
	for i := start + 1; i < l.Len()+1; i++ {
		l.hash[i-1] = l.hash[i]
	}
	delete(l.hash, l.Len()-1)
}

func (l list) pushTo(v interface{}, position int) *ListItem {
	prev := l.at(position - 1)
	next := l.at(position)
	item := &ListItem{
		Value:    v,
		Prev:     prev,
		Next:     next,
		position: position,
	}
	if prev != nil {
		prev.Next = item
	}
	if next != nil {
		next.Prev = item
	}
	l.incIndex(position)
	l.hash[position] = item
	return item
}
