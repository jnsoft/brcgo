package workers

import (
	"sync"

	"github.com/brcgo/src/domain"
	"github.com/brcgo/src/misc"
)

func ParserWorker(id int, lines <-chan string, parsedChans []chan domain.StringFloat, shardCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		data, _ := domain.ParseStringFloat(line)
		shard := misc.HashKey(data.Key) % shardCount
		parsedChans[shard] <- data
	}
}
