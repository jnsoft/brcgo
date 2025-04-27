package cache

type Item[K comparable, V any] struct {
	key      K
	value    V
	priority int
	index    int // Index in the heap (updated by heap)
}

type CacheMinHeap[K comparable, V any] []*Item[K, V]

func (pq CacheMinHeap[K, V]) Len() int { return len(pq) }

func (pq CacheMinHeap[K, V]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq CacheMinHeap[K, V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *CacheMinHeap[K, V]) Push(x any) {
	entry := x.(*Item[K, V])
	entry.index = len(*pq)
	*pq = append(*pq, entry)
}

func (pq *CacheMinHeap[K, V]) Pop() any {
	old := *pq
	n := len(old)
	entry := old[n-1]
	entry.index = -1 // Mark as removed
	*pq = old[:n-1]
	return entry
}
