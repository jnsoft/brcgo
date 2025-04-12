package workers

import (
	"math"
	"sync"

	"github.com/brcgo/src/models"
)

func LineWorker(id int, lines <-chan string, hashmap *map[string]models.StationData, mapMutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range lines {
		data := models.ParseLine(line)

		mapMutex.Lock()

		aggregated, exists := (*hashmap)[data.Key]
		if !exists {
			(*hashmap)[data.Key] = models.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			(*hashmap)[data.Key] = models.StationData{
				Min:   math.Min(data.Value, aggregated.Min),
				Max:   math.Max(data.Value, aggregated.Max),
				Sum:   data.Value + aggregated.Sum,
				Count: aggregated.Count + 1,
			}
		}

		mapMutex.Unlock()
	}
}
