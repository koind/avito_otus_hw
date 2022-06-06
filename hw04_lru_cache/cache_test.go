package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keyA Key = "aaa"
	keyB Key = "bbb"
	keyC Key = "ccc"
	keyD Key = "ddd"
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

	t.Run("purge logic", func(t *testing.T) {
		t.Run("purge first element", func(t *testing.T) {
			c := initCache(t)

			was := c.Set(keyD, 4)
			assert.False(t, was)

			a, ok := c.Get(keyA)
			assert.False(t, ok)
			assert.Nil(t, a)
		})

		t.Run("purge less used element", func(t *testing.T) {
			c := initCache(t)

			_ = c.Set(keyA, 4)
			_, _ = c.Get(keyC)

			was := c.Set(keyD, 4)
			assert.False(t, was)

			b, ok := c.Get(keyB)
			assert.False(t, ok)
			assert.Nil(t, b)
		})
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

func initCache(t *testing.T) Cache {
	t.Helper()

	c := NewCache(3)

	for i, key := range [...]Key{keyA, keyB, keyC} {
		was := c.Set(key, i)
		assert.False(t, was)
	}

	return c
}
