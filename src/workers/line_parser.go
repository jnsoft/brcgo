package workers

import (
	"fmt"
	"os"
)

func ParseLines[T any](in <-chan string, out chan<- T, parser func(string) (T, error)) {
	for line := range in {
		parsed, err := parser(line)
		if err == nil {
			out <- parsed
		} else {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		}
	}
	close(out)
}

func ParseLByteines[T any](in <-chan []byte, out chan<- T, parser func([]byte) (T, error)) {
	for line := range in {
		parsed, err := parser(line)
		if err == nil {
			out <- parsed
		} else {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		}
	}
	close(out)
}
