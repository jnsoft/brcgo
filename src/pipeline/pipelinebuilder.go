package pipeline

import (
	"fmt"
	"os"
)

// A source produces values of type T into a channel
type Source[T any] func(chan<- T) error

// A transform maps from T → U
type Transform[T, U any] func(<-chan T) <-chan U

// A side-effect stage that doesn’t change the type
type Effect[T any] func(T)

type PipelineBuilder[T any] struct {
	source Source[T]
	run    func(<-chan T)
}

func FromSource[T any](src func(chan<- T) error) PipelineBuilder[T] {
	return PipelineBuilder[T]{source: src}
}

func Then[T, U any](prev PipelineBuilder[T], transform func(<-chan T) <-chan U) PipelineBuilder[U] {
	return PipelineBuilder[U]{
		source: func(out chan<- U) error {
			in := make(chan T, 128)
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

func ThenDo[T any](prev PipelineBuilder[T], effect func(T)) PipelineBuilder[T] {
	return Then(prev, func(in <-chan T) <-chan T {
		out := make(chan T, 128)
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

func Finally[T any](pb PipelineBuilder[T], final func(T)) {
	pb.run = func(in <-chan T) {
		for val := range in {
			final(val)
		}
	}
}

func Run[T any](pb PipelineBuilder[T], finalize func()) {
	ch := make(chan T, 128)
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
