package workers

import (
	"bufio"
	"os"
)

func GetLines(filePath string, out chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		close(out)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
	return scanner.Err()
}

func GetByteLines(filePath string, out chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		close(out)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
	return scanner.Err()
}
