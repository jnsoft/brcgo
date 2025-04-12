package main

import (
	"fmt"

	"github.com/brcgo/src/util"
)

func main() {

	err := util.GenerateFile(5, "testfile.tmp")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Test file generated successfully.")
	}
}
