package workers

import (
	"sync"

	"github.com/brcgo/src/models"
	"github.com/brcgo/src/util"
)

func ParserWorker(id int, lines <-chan string, parsedChans []chan models.ParsedData, shardCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		data := models.ParseLine(line)
		shard := util.HashKey(data.Key) % shardCount
		parsedChans[shard] <- data
	}
}
