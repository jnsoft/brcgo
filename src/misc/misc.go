package misc

import (
	"math/rand"
	"sync"
	"time"
)

// Use a global random generator with a mutex for thread safety
var (
	globalRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randMutex  sync.Mutex
)

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
