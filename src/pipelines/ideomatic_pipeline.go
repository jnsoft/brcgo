package pipelines

import (
	"fmt"
	"os"

	"github.com/brcgo/src/workers"
)

func IdeomotaticPipeline[T, U any](fname string, parser func(string) (T, error), collector func(T), printer func(), verbose bool) {
	lines := make(chan string)
	parsed := make(chan T)

	go func() {
		if err := workers.GetLines(fname, lines); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		}
	}()

	go workers.ParseLines[T](lines, parsed, parser)

	workers.Collect(parsed, collector)

	printer()

}
