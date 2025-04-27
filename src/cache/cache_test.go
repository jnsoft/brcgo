package cache

import (
	"reflect"
	"sort"
	"testing"

	. "github.com/jnsoft/jngo/testhelper"
)

func TestSimpleCache(t *testing.T) {

	input_val := "value1"
	input_key := "key1"

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
		result := cache.GetMany(keys)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
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
