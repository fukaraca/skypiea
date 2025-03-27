package cache_test

import (
	"fmt"
	"github.com/fukaraca/skypiea/pkg/cache"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestStorage_SetGet(t *testing.T) {
	c := cache.New()

	c.Set("key1", "value1")
	c.Set("key2", 12345)
	c.Set("key3", []int{1, 2, 3})

	assert.Equal(t, "value1", c.Get("key1"))
	assert.Equal(t, 12345, c.Get("key2"))
	assert.Equal(t, []int{1, 2, 3}, c.Get("key3"))

	// Non-existing key
	assert.Nil(t, c.Get("unknown"))
}

func TestStorage_Del(t *testing.T) {
	c := cache.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	c.Del("key1")
	assert.Nil(t, c.Get("key1"))
	assert.Equal(t, "value2", c.Get("key2"))

	// Deleting non-existing key
	c.Del("key3") // Should not panic
	assert.Nil(t, c.Get("key3"))
}

func TestStorage_DeleteByPrefix(t *testing.T) {
	c := cache.New()

	c.Set("user_1", "Alice")
	c.Set("user_2", "Bob")
	c.Set("admin_1", "Charlie")
	c.Set("user_profile", "Dave")

	c.DeleteByPrefix("user_")

	assert.Nil(t, c.Get("user_1"))
	assert.Nil(t, c.Get("user_2"))
	assert.Nil(t, c.Get("user_profile"))
	assert.Equal(t, "Charlie", c.Get("admin_1"))
}

func TestStorage_Clear(t *testing.T) {
	c := cache.New()

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")

	c.Clear()

	assert.Nil(t, c.Get("key1"))
	assert.Nil(t, c.Get("key2"))
	assert.Nil(t, c.Get("key3"))
}

func TestStorage_Concurrency(t *testing.T) {
	c := cache.New()

	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(iterations * 2)

	for i := 0; i < iterations; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			c.Set(key, i)
		}(i)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			c.Get(key)
		}(i)
	}

	wg.Wait()

	// Confirm some keys exist
	assert.NotNil(t, c.Get("key1"))
	assert.NotNil(t, c.Get("key500"))
	assert.Equal(t, iterations, c.Len())
}
