package pipelines

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/brcgo/src/models"
	"github.com/brcgo/src/util"
	"github.com/brcgo/src/workers"
)

// Reading input and distributing it to a worker pool using goroutines and channels
// Wokers update the same map, sharing a mutex lock
func WorkerpoolPipeline(fname string, NO_OF_WORKERS int, verbose bool) {

	startTime := time.Now()

	lineChan := make(chan string)
	var wg sync.WaitGroup

	// Shared map and mutex
	resultMap := make(map[string]models.StationData)
	var mapMutex sync.Mutex

	// Start worker pool
	for i := 1; i <= NO_OF_WORKERS; i++ {
		wg.Add(1)
		go workers.LineWorker(i, lineChan, &resultMap, &mapMutex, &wg)
	}

	// Read file and send lines to channel
	err := util.ReadFileLines(fname, lineChan)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	close(lineChan) // No more lines coming
	wg.Wait()       // Wait for all workers to finish

	// Sort and print final results
	keys := make([]string, 0, len(resultMap))
	for k := range resultMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if verbose {
		fmt.Println("\n Final aggregated results:")
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, resultMap[k].String())
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, -1, len(resultMap))
}
