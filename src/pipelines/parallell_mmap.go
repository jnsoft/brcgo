package pipelines

import (
	"fmt"
	"sync"
	"time"

	"github.com/brcgo/src/domain"

	"golang.org/x/exp/mmap"
)

const (
	MMAP_BUFFER = 1024 * 1024 * 256
)

func ParallellMmap(fname string, max_cuncurrent int) (string, error) {
	startTime := time.Now()

	result := domain.NewByteResult()

	readTime := time.Now()

	elapsed := time.Since(readTime)
	fmt.Printf("\nRead time: %s", elapsed)

	var wg sync.WaitGroup

	for i := range max_cuncurrent {
		start := chunkStarts[i]
		end := chunkStarts[i+1]

		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			buf := make([]byte, end-start)
			_, err := reader.ReadAt(buf, int64(start))
			if err != nil {
				// handle error (optional: log or collect)
				return
			}
			ParseBuffer(buf, result)
		}(start, end)
	}

	waitTime := time.Now()
	wg.Wait()
	elapsed = time.Since(waitTime)
	fmt.Printf("\nWait time: %s", elapsed)

	sortTime := time.Now()
	res_str := result.GetSortedResults()
	elapsed = time.Since(sortTime)
	fmt.Printf("\nSort time: %s", elapsed)

	elapsed = time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, result.NoOfInputs(), result.NoOfStations())

	return res_str, nil
}

func splitFile(fname string, chunks int) ([]byte, error) {
	
	reader, err := mmap.Open(fname)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	
	size := reader.Len()
	chunkSize := size / chunks

	chunkStarts := make([]int, chunks+1)
	chunkStarts[0] = 0
	chunkStarts[chunks] = size
	for i := 1; i < chunks; i++ {
		pos := i * chunkSize
		for pos < size {
			b := make([]byte, 1)
			_, err := reader.ReadAt(b, int64(pos))
			if err != nil {
				break
			}
			if b[0] == ASCII_NEWLINE {
				pos++
				break
			}
			pos++
		}
		chunkStarts[i] = pos
	}
}
