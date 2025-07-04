package workers

func Collect[T any](in <-chan T, collector func(T)) {
	for item := range in {
		collector(item)
	}
}
