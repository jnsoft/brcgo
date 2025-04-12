package main

import (
	"fmt"

	"github.com/brcgo/src/misc"
	"github.com/brcgo/src/util"
)

func main() {

	for i := 0; i < 10; i++ {
		fmt.Println(misc.RandomInt(1, 10))
	}
	err := util.GenerateFile(5, "testfile.tmp")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Test file generated successfully.")
	}
}
