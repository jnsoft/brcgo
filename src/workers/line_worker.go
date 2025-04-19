package workers

import (
	"math"
	"sync"

	"github.com/brcgo/src/domain"
)

func LineWorker(id int, lines <-chan string, hashmap *map[string]domain.StationData, mapMutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range lines {
		data, _ := domain.ParseStringFloat(line)

		mapMutex.Lock()

		aggregated, exists := (*hashmap)[data.Key]
		if !exists {
			(*hashmap)[data.Key] = domain.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			(*hashmap)[data.Key] = domain.StationData{
				Min:   math.Min(data.Value, aggregated.Min),
				Max:   math.Max(data.Value, aggregated.Max),
				Sum:   data.Value + aggregated.Sum,
				Count: aggregated.Count + 1,
			}
		}

		mapMutex.Unlock()
	}
}
