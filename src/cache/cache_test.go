package cache

import (
	"reflect"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/jnsoft/jngo/testhelper"
)

const (
	input_val = "value1"
	input_key = "key1"
)

func TestSimpleCache(t *testing.T) {

	t.Run("Get, Set and Contains", func(t *testing.T) {

		cache := NewSimpleCache[string, string]()
		cache.Set(input_key, input_val)

		isFound := cache.Contains(input_key)
		AssertTrue(t, isFound)

		value, ok := cache.Get(input_key)
		if !ok || value != input_val {
			AssertEqual(t, value, input_val)
		}

		_, ok = cache.Get("nonexistent")
		AssertFalse(t, ok)
	})

	t.Run("Get and Set Many", func(t *testing.T) {
		cache := NewSimpleCache[string, string]()
		items := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		cache.SetMany(items)

		keys := []string{"key1", "key2", "key4"}
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		result, missing := cache.GetMany(keys)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		AssertEqual(t, len(missing), 1)
		AssertEqual(t, missing[0], "key4")
	})

	t.Run("Delete, Size and Clear", func(t *testing.T) {
		cache := NewSimpleCache[string, string]()
		cache.Set(input_key, input_val)
		cache.Delete(input_key)

		_, ok := cache.Get(input_key)
		AssertFalse(t, ok)

		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		size := cache.Size()
		AssertEqual(t, size, 2)

		cache.Clear()
		size = cache.Size()
		AssertEqual(t, size, 0)
	})

	t.Run("Keys and Values", func(t *testing.T) {
		cache := NewSimpleCache[string, string]()
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		keys := cache.Keys()
		expectedKeys := []string{"key1", "key2"}
		if !reflect.DeepEqual(keys, expectedKeys) {
			t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
		}

		values := cache.Values()
		sort.Strings(values)
		expectedValues := []string{"value1", "value2"}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Errorf("Expected values %v, got %v", expectedValues, values)
		}
	})
}

func TestLFUCache(t *testing.T) {

	t.Run("Get, Set and Contains", func(t *testing.T) {

		cache := NewLFUCache[string, int](2)
		cache.Set("A", 1)
		cache.Set("B", 2)

		isFound := cache.Contains("A")
		AssertTrue(t, isFound)

		if val, ok := cache.Get("A"); !ok || val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}

		cache.Set("C", 3)
		if _, ok := cache.Get("B"); ok {
			t.Errorf("Expected B to be evicted")
		}
	})

	t.Run("Get and Set Many", func(t *testing.T) {
		cache := NewLFUCache[string, int](3)

		cache.Set("X", 10)
		cache.Set("Y", 20)
		cache.Set("Z", 30)

		values, missing := cache.GetMany([]string{"X", "Y", "A"})

		if len(values) != 2 || values["X"] != 10 || values["Y"] != 20 {
			t.Errorf("Unexpected values: %v", values)
		}

		if len(missing) != 1 || missing[0] != "A" {
			t.Errorf("Unexpected missing keys: %v", missing)
		}

	})

	t.Run("SetOverCapacity", func(t *testing.T) {
		cache := NewLFUCache[string, int](1)

		cache.Set("A", 1)
		cache.Set("B", 2)

		if _, ok := cache.Get("A"); ok {
			t.Errorf("Expected A to be evicted")
		}
	})
}

func TestLRUCache(t *testing.T) {

	t.Run("Get, Set and Contains", func(t *testing.T) {

		cache := NewLRUCache[string, int](2)
		cache.Set("A", 1)
		cache.Set("B", 2)

		isFound := cache.Contains("A")
		AssertTrue(t, isFound)

		if val, ok := cache.Get("A"); !ok || val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}

		cache.Set("C", 3)
		if _, ok := cache.Get("B"); ok {
			t.Errorf("Expected B to be evicted since A was accessed last")
		}
	})

	t.Run("Get and Set Many", func(t *testing.T) {
		cache := NewLFUCache[string, int](3)

		cache.Set("X", 10)
		cache.Set("Y", 20)
		cache.Set("Z", 30)

		values, missing := cache.GetMany([]string{"X", "Y", "A"})

		if len(values) != 2 || values["X"] != 10 || values["Y"] != 20 {
			t.Errorf("Unexpected values: %v", values)
		}

		if len(missing) != 1 || missing[0] != "A" {
			t.Errorf("Unexpected missing keys: %v", missing)
		}

	})

	t.Run("SetOverCapacity", func(t *testing.T) {
		cache := NewLFUCache[string, int](1)

		cache.Set("A", 1)
		cache.Set("B", 2)

		if _, ok := cache.Get("A"); ok {
			t.Errorf("Expected A to be evicted")
		}
	})
}

func TestDLFUCache(t *testing.T) {
	t.Run("Get, Set and Contains", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("a", 1, time.Minute)

		test := cache.Contains("a")
		AssertTrue(t, test)

		test = cache.Contains("b")
		AssertFalse(t, test)

		v, ok := cache.Get("a")
		if !ok || v != 1 {
			t.Errorf("expected to get value 1 for key 'a', got %v (exists: %v)", v, ok)
		}
		_, ok = cache.Get("b")
		AssertFalse(t, ok)
	})

	t.Run("TestEviction", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("a", 1, time.Minute)
		cache.Set("b", 2, time.Minute)
		cache.Set("c", 3, time.Minute)

		size := cache.Size()
		AssertEqual(t, size, 2)
	})

	t.Run("TestFrequencyAffectsPriority", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("a", 1, time.Minute)
		cache.Set("b", 2, time.Minute)

		cache.Get("a")
		cache.Get("a")
		cache.Set("c", 3, time.Minute)

		_, ok := cache.Get("a")
		AssertTrue(t, ok)

		_, ok = cache.Get("b")
		AssertFalse(t, ok)
	})

	t.Run("TestExpiry", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("x", 42, 10*time.Millisecond)

		time.Sleep(20 * time.Millisecond)
		_, ok := cache.Get("x")
		AssertFalse(t, ok)
	})

	t.Run("Test Delete and Clear", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("x", 10, time.Minute)
		cache.Delete("x")

		_, ok := cache.Get("x")
		AssertFalse(t, ok)

		cache.Set("a", 1, time.Minute)
		cache.Set("b", 2, time.Minute)
		cache.Clear()

		size := cache.Size()
		AssertEqual(t, size, 0)
	})

	t.Run("TestKeysAndValues", func(t *testing.T) {
		cache := NewDLFUCache[string, int](2, 0.5)
		cache.Set("a", 1, time.Minute)
		cache.Set("b", 2, time.Minute)

		keys := cache.Keys()
		values := cache.Values()

		if len(keys) != 2 || len(values) != 2 {
			t.Errorf("expected 2 keys and 2 values, got %d keys and %d values", len(keys), len(values))
		}
	})

	t.Run("TestSetManyAndGetMany", func(t *testing.T) {
		cache := NewDLFUCache[string, int](3, 0.5)
		items := map[string]int{"a": 1, "b": 2, "c": 3}
		cache.SetMany(items, time.Minute)

		vals, missing := cache.GetMany([]string{"a", "b", "d"})

		if len(vals) != 2 {
			t.Errorf("expected 2 found keys, got %d", len(vals))
		}
		if len(missing) != 1 || missing[0] != "d" {
			t.Errorf("expected 'd' to be missing, got: %v", missing)
		}
	})
}

func BenchmarkSimpleCache(b *testing.B) {
	cache := NewSimpleCache[string, int]()
	keys := make([]string, 1000)
	items := make(map[string]int)
	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = i
	}

	b.ResetTimer()

	cache.SetMany(items)

	for i := 0; i < 1000; i++ {
		cache.Get(keys[i])
	}
}

func BenchmarkLFUCache(b *testing.B) {
	cache := NewLFUCache[string, int](1000)
	keys := make([]string, 1000)
	items := make(map[string]int)
	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = i
	}

	b.ResetTimer()

	cache.SetMany(items)

	for i := 0; i < 1000; i++ {
		cache.Get(keys[i])
	}

}

func BenchmarkLRUCache(b *testing.B) {
	cache := NewLRUCache[string, int](1000)
	keys := make([]string, 1000)
	items := make(map[string]int)
	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = i
	}

	b.ResetTimer()

	cache.SetMany(items)

	for i := 0; i < 1000; i++ {
		cache.Get(keys[i])
	}
}

func BenchmarkSimpleCacheParallel(b *testing.B) {
	cache := NewSimpleCache[string, int64]()
	keys := make([]string, 1000)
	items := make(map[string]int64)

	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = int64(i)
	}

	cache.SetMany(items)
	klen := int64(len(keys))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// set local state for goroutine here:
		var counter int64
		// execute benchmark iteration:
		for pb.Next() {
			index := atomic.AddInt64(&counter, 1) % klen
			cache.Get(keys[index])
		}
	})
}

func BenchmarkLFUCacheParallel(b *testing.B) {
	cache := NewLFUCache[string, int64](1000)
	keys := make([]string, 1000)
	items := make(map[string]int64)

	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = int64(i)
	}

	cache.SetMany(items)
	var counter int64
	klen := int64(len(keys))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := atomic.AddInt64(&counter, 1) % klen
			cache.Get(keys[index])
		}
	})
}

func BenchmarkLRUCacheParallel(b *testing.B) {
	cache := NewLRUCache[string, int64](1000)
	keys := make([]string, 1000)
	items := make(map[string]int64)

	for i := 0; i < 1000; i++ {
		keys[i] = "key" + strconv.Itoa(i)
		items[keys[i]] = int64(i)
	}

	cache.SetMany(items)
	var counter int64
	klen := int64(len(keys))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := atomic.AddInt64(&counter, 1) % klen
			cache.Get(keys[index])
		}
	})
}
