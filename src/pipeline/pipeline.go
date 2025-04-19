package pipeline

import (
	"fmt"
	"os"
	"sync"
)

// Stage is a processing step that transforms A → A (with optional side-effects).
type Stage[A any] func(<-chan A) <-chan A

// TransformStage is a pipeline step that converts A → B.
type TransformStage[A, B any] func(<-chan A) <-chan B

// MapStage maps items from type A to type B using a function.
func MapStage[A, B any](mapper func(A) (B, error)) TransformStage[A, B] {
	return func(in <-chan A) <-chan B {
		out := make(chan B, 128)
		go func() {
			defer close(out)
			for item := range in {
				if result, err := mapper(item); err == nil {
					out <- result
				} else {
					fmt.Fprintf(os.Stderr, "Map error: %v\n", err)
				}
			}
		}()
		return out
	}
}

// CollectStage applies a function to each item and passes it through (side-effect only).
func CollectStage[T any](fn func(T)) Stage[T] {
	return func(in <-chan T) <-chan T {
		out := make(chan T, 128)
		go func() {
			defer close(out)
			for item := range in {
				fn(item)
				out <- item
			}
		}()
		return out
	}
}

// TerminalStage applies a final side-effect and discards the data.
func TerminalStage[T any](fn func(T)) Stage[T] {
	return func(in <-chan T) <-chan T {
		done := make(chan T) // closed immediately; no output
		go func() {
			defer close(done)
			for item := range in {
				fn(item)
			}
		}()
		return done
	}
}

// ParallelMapStage processes items in parallel with worker count.
func ParallelMapStage[A, B any](workers int, mapper func(A) (B, error)) TransformStage[A, B] {
	return func(in <-chan A) <-chan B {
		out := make(chan B, 128)
		var wg sync.WaitGroup

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for item := range in {
					if result, err := mapper(item); err == nil {
						out <- result
					} else {
						fmt.Fprintf(os.Stderr, "Parallel map error: %v\n", err)
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
}

// ParallelCollectStage runs the collector function in parallel across N workers.
func ParallelCollectStage[T any](workers int, fn func(T)) Stage[T] {
	return func(in <-chan T) <-chan T {
		out := make(chan T, 128)
		var wg sync.WaitGroup

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for item := range in {
					fn(item)
					out <- item
				}
			}()
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
}

// Pipeline connects source → stages → final function.
// Each stage processes and passes data forward through buffered channels.
func Pipeline[A any](source func(chan<- A) error, stages []Stage[A], final func()) {
	src := make(chan A, 128)

	// Start source in a goroutine
	go func() {
		if err := source(src); err != nil {
			fmt.Fprintf(os.Stderr, "Error in source: %v\n", err)
		}
	}()

	var ch <-chan A = src // narrow type for read-only

	// Apply all stages
	for _, stage := range stages {
		ch = stage(ch)
	}

	// Drain final output
	for range ch {
		// if the last stage is TerminalStage, this is empty
		// Final result is consumed or discarded
	}

	// Optional final hook
	if final != nil {
		final()
	}
}
