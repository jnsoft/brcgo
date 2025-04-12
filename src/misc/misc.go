package misc

import (
	"math/rand"
	"time"
)

func RandomInt(min, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(max-min+1) + min
}

func GetRandomName(lg int) string {
	const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lower = "abcdefghijklmnopqrstuvwxyz"

	var res = make([]byte, lg)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	res[0] = upper[r.Intn(len(upper))]
	for i := 1; i < lg; i++ {
		res[i] = lower[r.Intn(len(lower))]
	}

	return string(res)
}
