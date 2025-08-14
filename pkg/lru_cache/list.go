package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Length int
	First  *ListItem
	Last   *ListItem
}

func NewList() List {
	return &list{
		Length: 0,
		First:  nil,
		Last:   nil,
	}
}

func NewListItem(v any) *ListItem {
	return &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v any) *ListItem {
	newLI := NewListItem(v)
	if l.Length == 0 {
		l.First = newLI
		l.Last = newLI
	} else {
		newLI.Next = l.First

		l.First.Prev = newLI

		l.First = newLI
	}
	l.Length++
	return newLI
}

func (l *list) PushBack(v any) *ListItem {
	newLI := NewListItem(v)
	if l.Length == 0 {
		l.First = newLI
		l.Last = newLI
	} else {
		newLI.Prev = l.Last

		l.Last.Next = newLI

		l.Last = newLI
	}
	l.Length++
	return newLI
}

func (l *list) Remove(i *ListItem) {
	if l.Length == 0 {
		return
	}
	if i.Next == nil && i.Prev == nil {
		l.Length = 0
		l.Last = nil
		l.First = nil

		return
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.Last = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.First = i.Next
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		l.Last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Prev.Next = i.Next
	i.Prev = nil
	i.Next = l.First
	l.First.Prev = i
	l.First = i
}
