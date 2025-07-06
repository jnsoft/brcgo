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

	readTime := time.Now()
	data, err := splitFile(fname, max_cuncurrent)
	if err != nil {
		return "", err
	}
	elapsed := time.Since(readTime)
	fmt.Printf("\nRead time: %s", elapsed)

	results := make([]*domain.ByteResult2, max_cuncurrent)
	for i := range results {
		results[i] = domain.NewByteResult2()
	}

	var wg sync.WaitGroup

	for i := range max_cuncurrent {

		wg.Add(1)
		go func(buf []byte, result *domain.ByteResult2) {
			defer wg.Done()
			parseBuffer(buf, result)
		}(data[i], results[i])
	}

	waitTime := time.Now()
	wg.Wait()
	elapsed = time.Since(waitTime)
	fmt.Printf("\nWait time: %s", elapsed)

	mergeTime := time.Now()
	for i := 1; i < len(results); i++ {
		results[0].Merge(results[i])
	}
	elapsed = time.Since(mergeTime)
	fmt.Printf("\nMerge time: %s", elapsed)

	sortTime := time.Now()
	res_str := results[0].GetSortedResults()
	elapsed = time.Since(sortTime)
	fmt.Printf("\nSort time: %s", elapsed)

	elapsed = time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, results[0].NoOfInputs(), results[0].NoOfStations())

	return res_str, nil
}

func splitFile(fname string, chunks int) ([][]byte, error) {

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

	b := make([]byte, 1)
	for i := 1; i < chunks; i++ {
		pos := i * chunkSize
		for pos < size {
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

	result := make([][]byte, chunks)

	for i := 0; i < chunks; i++ {
		start := chunkStarts[i]
		end := chunkStarts[i+1]
		length := end - start

		buf := make([]byte, length)
		_, err := reader.ReadAt(buf, int64(start))
		if err != nil {
			return nil, err
		}
		result[i] = buf
	}

	return result, nil
}

func parseBuffer(parseBuffer []byte, result *domain.ByteResult2) {
	lineStartIdx := 0
	for i := 0; i < len(parseBuffer); i++ {
		if parseBuffer[i] == ASCII_NEWLINE {
			lineEndIdx := i
			// Handle \r\n (Windows line endings)
			if lineEndIdx > lineStartIdx && parseBuffer[lineEndIdx-1] == '\r' {
				lineEndIdx--
			}
			line := parseBuffer[lineStartIdx:lineEndIdx]
			if len(line) > 0 {
				reading := domain.NewByteStationReadingFromBytes(line)
				result.Add(reading)
			}
			lineStartIdx = i + 1
		}
	}
}
