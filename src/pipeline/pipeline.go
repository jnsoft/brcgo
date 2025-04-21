package pipeline

import (
	"fmt"
	"os"
	"sync"
)

const (
	BUFFER_SIZE = 128
)

/*
Manual pipelines with Pipeline(...)
Builder-style pipelines with FromSource(...)
*/

type PipelineBuilder[T any] struct {
	source Source[T]
	run    func(<-chan T)
}

// --- Pipeline Builder (Fluent Interface) ---

// FromSource starts a pipeline from a source.
func FromSource[T any](src func(chan<- T) error) PipelineBuilder[T] {
	return PipelineBuilder[T]{source: src}
}

// Then adds a transformation to the pipeline.
func Then[T, U any](prev PipelineBuilder[T], transform func(<-chan T) <-chan U) PipelineBuilder[U] {
	return PipelineBuilder[U]{
		source: func(out chan<- U) error {
			in := make(chan T, BUFFER_SIZE)
			outChan := transform(in)

			errCh := make(chan error, 1)
			go func() {
				err := prev.source(in)
				//close(in)
				errCh <- err
			}()

			for val := range outChan {
				out <- val
			}
			close(out)

			return <-errCh
		},
	}
}

// ThenDo adds a side-effect (non-transforming) stage.
func ThenDo[T any](prev PipelineBuilder[T], effect func(T)) PipelineBuilder[T] {
	return Then(prev, func(in <-chan T) <-chan T {
		out := make(chan T, BUFFER_SIZE)
		go func() {
			defer close(out)
			for val := range in {
				effect(val)
				out <- val
			}
		}()
		return out
	})
}

// Finally adds a final consumer for the data.
func Finally[T any](pb PipelineBuilder[T], final func(T)) PipelineBuilder[T] {
	pb.run = func(in <-chan T) {
		for val := range in {
			final(val)
		}
	}
	return pb
}

// Run starts the pipeline and calls finalize when complete.
func Run[T any](pb PipelineBuilder[T], finalize func()) {
	ch := make(chan T, BUFFER_SIZE)
	go func() {
		if err := pb.source(ch); err != nil {
			fmt.Fprintf(os.Stderr, "Pipeline error: %v\n", err)
		}
	}()
	if pb.run != nil {
		pb.run(ch)
	} else {
		for range ch {
		}
	}
	if finalize != nil {
		finalize()
	}
}

// --- Low-level pipeline runner ---
// Pipeline connects source → stages → final function.
// Each stage processes and passes data forward through buffered channels.
func Pipeline[A any](source func(chan<- A) error, stages []Stage[A], final func()) {
	src := make(chan A, BUFFER_SIZE)

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

// --- Core types ---
// Source produces values of type T into a channel
type Source[T any] func(chan<- T) error

// Stage is a processing step that transforms T → T (with optional side-effects).
type Stage[T any] func(<-chan T) <-chan T

// Transform is a pipeline step that maps from T → U
type Transform[T, U any] func(<-chan T) <-chan U

// A side-effect stage that doesn’t change the type
type Effect[T any] func(T)

// --- Functional helpers ---

// MapStage maps items from A to B using a function, emits only successful results.
func MapStage[A, B any](mapper func(A) (B, error)) Transform[A, B] {
	return func(in <-chan A) <-chan B {
		out := make(chan B, BUFFER_SIZE)
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
		out := make(chan T, BUFFER_SIZE)
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
func ParallelMapStage[A, B any](workers int, mapper func(A) (B, error)) Transform[A, B] {
	return func(in <-chan A) <-chan B {
		out := make(chan B, BUFFER_SIZE)
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

// ParallelCollectStage runs side-effects in parallel.
func ParallelCollectStage[T any](workers int, fn func(T)) Stage[T] {
	return func(in <-chan T) <-chan T {
		out := make(chan T, BUFFER_SIZE)
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

func ParallelDoStage[T any](workers int, fn func(T)) func(<-chan T) <-chan T {
	return func(in <-chan T) <-chan T {
		out := make(chan T)
		var wg sync.WaitGroup

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for val := range in {
					fn(val)
					out <- val // optional: forward data to next stage
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
