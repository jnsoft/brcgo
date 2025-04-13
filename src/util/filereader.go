package util

import (
	"bufio"
	"os"
)

func ReadFileLines(fname string, out chan<- string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}

	return scanner.Err()
}

func ReadFileBytes(fname string, out chan<- []byte) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Bytes()
	}

	return scanner.Err()
}
