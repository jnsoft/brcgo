package pipelines

import (
	"fmt"
	"time"

	"golang.org/x/exp/mmap"
)

const (
	MMAP_BUFFER_SIZE = 1024 * 16
)

func SimpleMmap(fname string) error {
	startTime := time.Now()

	reader, err := mmap.Open(fname)
	if err != nil {
		return err
	}
	defer reader.Close()

	size := reader.Len()
	cnt := 0
	offset := 0
	buf := make([]byte, MMAP_BUFFER_SIZE)

	for offset < size {
		toRead := min(size-offset, MMAP_BUFFER_SIZE)
		n, err := reader.ReadAt(buf[:toRead], int64(offset))
		if err != nil && n == 0 {
			break
		}
		for _, b := range buf[:n] {
			if b == ASCII_NEWLINE {
				cnt++
			}
		}
		offset += n
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, cnt, 0)

	return nil
}
