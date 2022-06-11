package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Get(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type value interface{}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items map[value]*ListItem
	first *ListItem
	last  *ListItem
}

// Creates a new list.
func NewList() List {
	return &list{
		items: make(map[value]*ListItem),
	}
}

// Returns the number of elements in the list.
func (l *list) Len() int {
	return len(l.items)
}

// Returns an element by value.
func (l *list) Get(v interface{}) *ListItem {
	val, ok := v.(value)
	if !ok {
		return nil
	}

	i, ok := l.items[val]
	if !ok {
		return nil
	}

	return i
}

// Returns the first element of the list.
func (l *list) Front() *ListItem {
	return l.first
}

// Returns the last element of the list.
func (l *list) Back() *ListItem {
	return l.last
}

// Adds a new element to the beginning of the list.
func (l *list) PushFront(v interface{}) *ListItem {
	if v == nil {
		return nil
	}

	i := &ListItem{Value: v}

	if l.Len() > 0 {
		l.first.Prev = i
		i.Next = l.first
	}

	l.first = i
	if l.Len() == 0 {
		l.last = i
	}

	l.items[i.Value] = i

	return i
}

// Adds a new element to the end of the list.
func (l *list) PushBack(v interface{}) *ListItem {
	if v == nil {
		return nil
	}

	i := &ListItem{Value: v}

	if l.Len() > 0 {
		l.last.Next = i
		i.Prev = l.last
	}

	l.last = i

	if l.Len() == 0 {
		l.first = i
	}

	l.items[i.Value] = i

	return i
}

// Removes an element from the list.
func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	delete(l.items, i.Value)

	if i == l.first {
		l.first.Next.Prev = nil
		l.first = l.first.Next
		return
	}

	if i == l.last {
		l.last.Prev.Next = nil
		l.last = l.last.Prev
		return
	}

	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next
}

// Moves the element to the beginning.
func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}
