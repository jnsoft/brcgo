package util

import "hash/fnv"

// Fowler-Noll-Vo hash (FNV-1a) algorithm
func HashKey(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}
