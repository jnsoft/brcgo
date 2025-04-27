package cache

import (
	"reflect"
	"sort"
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

		time.Sleep(1 * time.Second) // Ensure some time has passed for cache eviction

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
