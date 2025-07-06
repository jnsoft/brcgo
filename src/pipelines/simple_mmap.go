package pipelines

import (
	"fmt"
	"time"

	"github.com/brcgo/src/domain"

	"golang.org/x/exp/mmap"
)

func SimpleMmap(fname string) (string, error) {
	startTime := time.Now()

	reader, err := mmap.Open(fname)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	result := domain.NewByteResult()

	size := reader.Len()
	buf := make([]byte, size)

	n, err := reader.ReadAt(buf, 0)
	if err != nil && n == 0 {
		return "", err
	}

	lineStartIdx := 0
	for i := range size {
		if buf[i] == ASCII_NEWLINE {
			lineEndIdx := i

			//Only needed if we want to handle \r\n line endings
			//if lineEndIdx > lineStartIdx && buf[lineEndIdx-1] == '\r' {
			//	lineEndIdx--
			//}

			reading := domain.NewByteStationReadingFromBytes(buf[lineStartIdx:lineEndIdx])
			result.Add(reading)

			lineStartIdx = i + 1
		}
	}

	res_str := result.GetSortedResults()

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, result.NoOfInputs(), result.NoOfStations())

	return res_str, nil
}
