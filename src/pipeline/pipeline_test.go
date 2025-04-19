package pipeline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brcgo/src/misc"
	"github.com/brcgo/src/workers"
)

const FNAME = "delteme.txt"

type SampleData struct {
	Key   string
	Value float64
}

func generateTestFile(path string, lines int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < lines; i++ {
		_, _ = fmt.Fprintf(f, "sensor-%d;%.2f\n", misc.RandomInt(0, 999), misc.Random()*100)
	}
	return nil
}

func parseLine(line string) (SampleData, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		return SampleData{}, fmt.Errorf("invalid line: %s", line)
	}
	val, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return SampleData{}, err
	}
	return SampleData{Key: parts[0], Value: val}, nil
}

func TestPipelineBuilder(t *testing.T) {

	// Generate test input

	if err := generateTestFile(FNAME, 1_000); err != nil {
		panic(err)
	}

	counter := make(map[string]int)

	start := time.Now()

	pb := FromSource(func(out chan<- string) error {
		return workers.GetLines(FNAME, out)
	})

	pb2 := Then(pb, ParallelMapStage[string, SampleData](8, parseLine))

	pb3 := ThenDo(pb2, func(data SampleData) {
		counter[data.Key]++
	})

	Finally(pb3, func(data SampleData) {
		// drain if needed
	})

	Run(pb3, func() {
		fmt.Printf("Processed %d keys in %v\n", len(counter), time.Since(start))
	})

	err := os.Remove(FNAME)
	if err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
	} else {
		fmt.Println("File deleted successfully")
	}
}

func TestPipeline(t *testing.T) {

}
