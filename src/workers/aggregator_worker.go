package workers

import (
	"math"
	"sync"

	"github.com/brcgo/src/models"
)

type AggregatorStats struct {
	ItemsProcessed int
	UniqueKeys     int
}

type AggregatorResult struct {
	ID    int
	Data  map[string]models.StationData
	Stats AggregatorStats
}

func AggregatorWorker(id int, input <-chan models.ParsedData, out chan<- AggregatorResult, wg *sync.WaitGroup) {
	defer wg.Done()

	hashmap := make(map[string]models.StationData)
	var stats AggregatorStats

	for data := range input {
		aggregated, exists := hashmap[data.Key]
		if !exists {
			hashmap[data.Key] = models.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			hashmap[data.Key] = models.StationData{
				Min:   math.Min(data.Value, aggregated.Min),
				Max:   math.Max(data.Value, aggregated.Max),
				Sum:   data.Value + aggregated.Sum,
				Count: aggregated.Count + 1,
			}
		}
		stats.ItemsProcessed++
	}

	stats.UniqueKeys = len(hashmap)

	out <- AggregatorResult{
		ID:    id,
		Data:  hashmap,
		Stats: stats,
	}
}
