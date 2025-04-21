package parallell

import "sync"

func LimitedParallelRun[T any](items []T, maxConcurrency int, fn func(T)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	for _, item := range items {
		wg.Add(1)
		sem <- struct{}{}
		go func(item T) {
			defer wg.Done()
			fn(item)
			<-sem
		}(item)
	}

	wg.Wait()
}
