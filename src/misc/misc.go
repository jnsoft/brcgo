package misc

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"runtime/pprof"
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

// view profile: go tool pprof -http 127.0.0.1:8080 ./cpu_profile.prof
func ProfileFunction(label, prof_file string, f func() (any, error)) {
	pfile, err := os.Create(prof_file)
	if err != nil {
		panic(err)
	}
	defer pfile.Close()

	if err := pprof.StartCPUProfile(pfile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	start := time.Now()
	result, err := f()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%s: %v (%s)\n", label, result, elapsed)

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

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Fowler-Noll-Vo hash (FNV-1a) algorithm
func HashKey(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}
