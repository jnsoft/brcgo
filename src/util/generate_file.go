package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/jnsoft/jngo/randx"
)

func GenerateFile(size, no_of_locations int, fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 1<<20) // 1MB buffer

	locations := get_locations(no_of_locations)

	for i := 0; i < size; i++ {
		rowData := fmt.Sprintf("%s;%s\n", locations[randx.Int(0, no_of_locations-1)], get_temp())
		_, err := writer.WriteString(rowData)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	fmt.Println("File created successfully:", fname)
	return nil
}

func get_temp() string {
	t := randx.Int(-999, 999)
	f := float64(t) / 10.0
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func get_locations(size int) []string {
	var res = make([]string, size)
	for i := 0; i < size; i++ {
		res[i] = randx.GetRandomName(randx.Int(3, 14))
	}
	return res
}

/*

Yellowknife;16.0
Entebbe;32.9

*/
