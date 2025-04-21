package pipeline

import (
	"bufio"
	"os"
)

func FileLineSource(path string) Source[string] {
	return func(out chan<- string) error {
		defer close(out)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		return scanner.Err()
	}
}

func SliceSource[T any](items []T) Source[T] {
	return func(out chan<- T) error {
		defer close(out)
		for _, item := range items {
			out <- item
		}
		return nil
	}
}
