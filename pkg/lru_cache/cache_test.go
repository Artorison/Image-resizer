package lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
	c.Clear()
	_ = t
}

func TestCustomCache(t *testing.T) {
	t.Run("max buffer", func(t *testing.T) {
		c := NewCache(3)

		_ = c.Set("aaa", 111)
		_ = c.Set("bbb", 222)
		_ = c.Set("ccc", 333)
		_ = c.Set("ddd", 444)
		_ = c.Set("fff", 555)

		val, ok := c.Get("aaa")

		assert.False(t, ok)
		assert.Nil(t, val)
		c.Clear()
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		_ = c.Set("aaa", 111)
		_ = c.Set("bbb", 222)
		_ = c.Set("ccc", 333)

		_ = c.Set("bbb", 333)

		val, ok := c.Get("aaa")
		assert.True(t, ok)
		assert.Equal(t, 111, val)

		_ = c.Set("aaa", 555)
		_, _ = c.Get("bbb")

		val, ok = c.Get("aaa")
		assert.True(t, ok)
		assert.Equal(t, 555, val)

		_ = c.Set("ccc", 555)
		val, ok = c.Get("ccc")
		assert.True(t, ok)
		assert.Equal(t, 555, val)

		_ = c.Set("fff", 999)

		val, ok = c.Get("bbb")
		assert.False(t, ok)
		assert.Nil(t, val)

		c.Clear()
	})

	t.Run("many purges", func(t *testing.T) {
		c := NewCache(10000)
		wg := &sync.WaitGroup{}

		for i := 1; i <= 10000; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c.Set(strconv.Itoa(i), i)
			}(i)
		}

		wg.Wait()

		val, ok := c.Get("900")
		assert.True(t, ok)
		assert.Equal(t, 900, val)

		val, ok = c.Get("1")
		assert.True(t, ok)
		assert.Equal(t, 1, val)

		val, ok = c.Get("9999")
		assert.True(t, ok)
		assert.Equal(t, 9999, val)

		for i := 10000; i <= 20000; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c.Set(strconv.Itoa(i), i)
			}(i)
		}

		wg.Wait()
		val1, ok1 := c.Get("900")
		assert.False(t, ok1)
		assert.Nil(t, val1)

		val1, ok1 = c.Get("1")
		assert.False(t, ok1)
		assert.Nil(t, val1)

		val1, ok1 = c.Get("9999")
		assert.False(t, ok1)
		assert.Nil(t, val1)

		c.Clear()
	})
}
