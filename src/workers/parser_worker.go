package workers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/brcgo/src/models"
	"github.com/brcgo/src/util"
)

func ParserWorker(id int, lines <-chan string, parsedChans []chan models.ParsedData, shardCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			fmt.Printf("Parser %d: Invalid line: %s\n", id, line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])

		value, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			fmt.Printf("Parser %d: Failed to parse float: %s\n", id, line)
			continue
		}

		data := models.ParsedData{Key: key, Value: value}
		shard := util.HashKey(key) % shardCount
		parsedChans[shard] <- data
	}
}
