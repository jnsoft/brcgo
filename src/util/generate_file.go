package util

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brcgo/src/misc"
)

func GenerateFile(size, no_of_locations int, fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	locations := get_locations(no_of_locations)

	for range size {
		rowData := fmt.Sprintf("%s;%s", locations[misc.RandomInt(0, no_of_locations-1)], get_temp())
		_, err := file.WriteString(rowData + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	fmt.Println("File created successfully:", fname)
	return nil
}

func get_temp() string {
	t := misc.RandomInt(-999, 999)
	f := float64(t) / 10.0
	return strconv.FormatFloat(f, 'f', 1, 64)
}

func get_locations(size int) []string {
	var res = make([]string, size)
	for i := 0; i < size; i++ {
		res[i] = misc.GetRandomName(misc.RandomInt(3, 14))
	}
	return res
}

/*

Yellowknife;16.0
Entebbe;32.9

*/
