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
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length      int
	first, last *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.first,
	}
	l.first = item
	l.length++
	if l.length > 1 {
		item.Next.Prev = item
	} else {
		l.last = l.first
	}
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.last,
	}
	l.last = item
	l.length++
	if l.length > 1 {
		item.Prev.Next = item
	} else {
		l.last = l.first
	}
	return item
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		break
	case i == l.last:
		l.last = i.Prev
		i.Prev.Next = nil
		i = nil
	case i == l.first:
		l.first = i.Next
		i.Next.Prev = nil
	}
	l.length--
	i = nil
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
