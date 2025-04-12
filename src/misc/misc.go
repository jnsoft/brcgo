package misc

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Use a global random generator with a mutex for thread safety
var (
	globalRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randMutex  sync.Mutex
)

func TimeFunction(label string, f func() (any, error)) {
	start := time.Now()
	result, err := f()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%s: %v (%s)\n", label, result, elapsed)
}

func TimeFunctionNoResult(label string, f func() error) {
	start := time.Now()
	err := f()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%s completed in %s\n", label, elapsed)
}

func TimeFunctionVoid(label string, f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	fmt.Printf("%s completed in %s\n", label, elapsed)
}

func RandomInt(min, max int) int {
	randMutex.Lock()
	defer randMutex.Unlock()

	return globalRand.Intn(max-min+1) + min
}

func GetRandomName(lg int) string {
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lower = "abcdefghijklmnopqrstuvwxyz"

	var res = make([]byte, lg)

	randMutex.Lock()
	defer randMutex.Unlock()

	res[0] = upper[globalRand.Intn(len(upper))]
	for i := 1; i < lg; i++ {
		res[i] = lower[globalRand.Intn(len(lower))]
	}

	return string(res)
}
