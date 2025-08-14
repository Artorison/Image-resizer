package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestCustom(t *testing.T) {
	t.Run("move head", func(t *testing.T) {
		l := NewList()

		l.PushFront(100)
		l.PushBack(200)
		require.Equal(t, 2, l.Len())

		li := NewListItem(100)
		l.MoveToFront(li)
		require.Equal(t, 100, l.Front().Value)
	})

	t.Run("move nil", func(t *testing.T) {
		l := NewList()

		l.PushFront(100)
		l.PushBack(200)
		require.Equal(t, 2, l.Len())

		l.MoveToFront(&ListItem{})
		require.Equal(t, 2, l.Len())
	})
}

func TestLength(t *testing.T) {
	t.Run("many add", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 100; i++ {
			l.PushFront(i)
			l.PushBack(i)
		}
		require.Equal(t, 200, l.Len())
	})
	t.Run("many delete", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 200; i++ {
			l.PushFront(i)
		}
		for i := 0; i < 100; i++ {
			l.Remove(l.Front())
		}
		require.Equal(t, 100, l.Len())
	})

	t.Run("delete all", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 100; i++ {
			l.PushBack(i)
		}
		for i := 0; i < 100; i++ {
			l.Remove(l.Front())
		}
		require.Equal(t, 0, l.Len())
	})

	t.Run("delete more then we have", func(t *testing.T) {
		l := NewList()
		for i := 0; i < 20; i++ {
			l.PushFront(i)
		}
		for i := 0; i < 100; i++ {
			l.Remove(l.Front())
		}
		require.Equal(t, 0, l.Len())
	})
}

func TestMoveToFront(t *testing.T) {
	t.Run("many add and MoveToFront", func(t *testing.T) {
		l := NewList()
		for i := 1; i <= 100; i++ {
			l.PushFront(i)
		}

		l.MoveToFront(l.Back())
		l.MoveToFront(l.Back())
		l.MoveToFront(l.Back())

		elems := make([]int, 0, 4)
		for item, i := l.Front(), 0; i < 4; item = item.Next {
			elems = append(elems, item.Value.(int))
			i++
		}
		require.Equal(t, []int{3, 2, 1, 100}, elems)
	})
}
