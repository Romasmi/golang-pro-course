package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

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
		shouldPanic(t, func() {
			_ = NewCache(0)
		})

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

	t.Run("purge logic - 1 capacity", func(t *testing.T) {
		c := NewCache(1)
		require.False(t, c.Set("a", 61))

		v, exist := c.Get("a")
		require.True(t, exist)
		require.Equal(t, v, 61)

		c.Set("b", 62)
		v, exist = c.Get("b")
		require.True(t, exist)
		require.Equal(t, v, 62)

		require.False(t, c.Set("a", 61))
		v, exist = c.Get("a")
		require.True(t, exist)
		require.Equal(t, v, 61)
	})

	t.Run("purge logic - N capacity", func(t *testing.T) {
		c := NewCache(3)
		require.False(t, c.Set("a", 61))

		v, exist := c.Get("a")
		require.True(t, exist)
		require.Equal(t, v, 61)

		c.Set("b", 62)
		v, exist = c.Get("b")
		require.True(t, exist)
		require.Equal(t, v, 62)

		c.Set("c", 63)
		v, exist = c.Get("c")
		require.True(t, exist)
		require.Equal(t, v, 63)

		// exceed capacity - "a" should be removed, "b" still should exist
		require.False(t, c.Set("d", 64))
		v, exist = c.Get("d")
		require.True(t, exist)
		require.Equal(t, v, 64)

		// b still exists
		v, exist = c.Get("b")
		require.True(t, exist)
		require.Equal(t, v, 62)

		// b purged
		require.False(t, c.Set("a", 61))
		v, exist = c.Get("a")
		require.True(t, exist)
		require.Equal(t, v, 61)

		v, exist = c.Get("b")
		require.False(t, exist)
		require.Equal(t, v, nil)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic as expected")
		}
	}()
	f()
}
