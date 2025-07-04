package pipelines

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/brcgo/src/domain"
)

const BUFFER_SIZE = 1024 * 1024
const ASCII_NEWLINE = '\n'

func NaiveBytes(fname string, MAX_CONCURRENT int) (string, error) {

	startTime := time.Now()

	file, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer file.Close()

	result := domain.NewByteResult()
	buffer := make([]byte, BUFFER_SIZE)
	var leftover []byte
	var wg sync.WaitGroup
	sem := make(chan struct{}, MAX_CONCURRENT)

	for {
		bytesRead, err := file.Read(buffer)
		if bytesRead == 0 && err != nil {
			break
		}

		// combine leftover with current buffer
		combined := append(leftover, buffer[:bytesRead]...)

		// Find last newline
		lastNewline := -1
		for i := len(combined) - 1; i >= 0; i-- {
			if combined[i] == ASCII_NEWLINE {
				lastNewline = i
				break
			}
		}

		if lastNewline == -1 {
			leftover = combined
			if err != nil {
				break
			}
			continue
		}

		leftover = nil
		if lastNewline+1 < len(combined) {
			leftover = combined[lastNewline+1:]
		}

		parseBuffer := make([]byte, lastNewline+1)
		copy(parseBuffer, combined[:lastNewline+1])

		wg.Add(1)
		sem <- struct{}{} // Acquire a semaphore slot
		go func(buf []byte) {
			defer wg.Done()
			defer func() { <-sem }() // Release the semaphore slot
			ParseBuffer(buf, result)
		}(parseBuffer)

		if err != nil {
			break
		}
	}
	if len(leftover) > 0 {
		reading := domain.NewByteStationReadingFromBytes(leftover)
		result.Add(reading)
	}

	wg.Wait()
	res_str := result.GetSortedResults()

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, result.NoOfInputs(), result.NoOfStations())

	return res_str, nil
}

func ParseBuffer(parseBuffer []byte, result *domain.ByteResult) {
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
