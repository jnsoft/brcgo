package util

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brcgo/src/misc"
)

func GenerateFile(size int, fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	rowData := fmt.Sprintf("%s;%s", misc.GetRandomName(misc.RandomInt(3, 14)), get_temp())
	for i := 0; i < size; i++ {
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

/*

Yellowknife;16.0
Entebbe;32.9

*/
