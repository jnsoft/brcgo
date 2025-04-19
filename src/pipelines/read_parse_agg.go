package pipelines

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/brcgo/src/domain"
	"github.com/brcgo/src/workers"
)

func ReadParseAggregatePipeline(fname string, NO_OF_PARSER_WORKERS, NO_OF_AGGREGATOR_WORKERS int, verbose bool) {

	startTime := time.Now()

	if verbose {
		fmt.Println("Setup plumbing...")
	}

	lineChan := make(chan string)
	parsedChans := make([]chan domain.StringFloat, NO_OF_AGGREGATOR_WORKERS)
	resultChan := make(chan workers.AggregatorResult, NO_OF_AGGREGATOR_WORKERS)

	// Create aggregator channels
	for i := range parsedChans {
		parsedChans[i] = make(chan domain.StringFloat, 100)
	}

	// Start aggregators
	var wgAggregators sync.WaitGroup
	for i := 0; i < NO_OF_AGGREGATOR_WORKERS; i++ {
		wgAggregators.Add(1)
		go workers.AggregatorWorker(i, parsedChans[i], resultChan, &wgAggregators)
	}

	// Start parsers
	var wgParsers sync.WaitGroup
	for i := 0; i < NO_OF_PARSER_WORKERS; i++ {
		wgParsers.Add(1)
		go workers.ParserWorker(i, lineChan, parsedChans, NO_OF_AGGREGATOR_WORKERS, &wgParsers)
	}

	if verbose {
		fmt.Println("Starting pipeline...")
	}

	// Reader
	err := workers.GetLines(fname, lineChan)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	close(lineChan)

	wgParsers.Wait()
	for _, ch := range parsedChans {
		close(ch)
	}

	wgAggregators.Wait()
	close(resultChan)

	// Combine results
	finalMap := make(map[string]domain.StationData)
	var totalStats workers.AggregatorStats
	for res := range resultChan {
		if verbose {
			fmt.Printf("Aggregator %d stats: %d items, %d keys\n",
				res.ID, res.Stats.ItemsProcessed, res.Stats.UniqueKeys)
		}

		for k, v := range res.Data {
			value, exists := finalMap[k]
			if !exists {
				finalMap[k] = domain.StationData{
					Min:   v.Min,
					Max:   v.Max,
					Sum:   v.Sum,
					Count: v.Count,
				}
			} else {
				finalMap[k] = domain.StationData{
					Min:   math.Min(v.Min, value.Min),
					Max:   math.Max(v.Max, value.Max),
					Sum:   value.Sum + v.Sum,
					Count: value.Count + v.Count,
				}
			}

		}

		totalStats.ItemsProcessed += res.Stats.ItemsProcessed
		totalStats.UniqueKeys += res.Stats.UniqueKeys // may have overlap
	}

	// Sort and print final results
	keys := make([]string, 0, len(finalMap))
	for k := range finalMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if verbose {
		fmt.Println("\n Final aggregated results:")
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, finalMap[k].String())
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, approx. %d unique keys\n",
		elapsed, totalStats.ItemsProcessed, len(finalMap))
}
