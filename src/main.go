package main

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brcgo/src/util"
)

const (
	NO_OF_PARSER_WORKERS     = 10
	NO_OF_AGGREGATOR_WORKERS = 3
	ERROR                    = "❌ Error reading file"
	WARNING                  = "⚠️ Warning"
	DONE                     = "✅ Done"
)

type ParsedData struct {
	Key   string
	Value float64
}

type AggregatorResult struct {
	ID    int
	Data  map[string]float64
	Stats AggregatorStats
}

type AggregatorStats struct {
	LinesProcessed int
	UniqueKeys     int
	TotalValue     float64
}

func main() {
	fname := "testfile_100.tmp"

	if false {
		util.GenerateFile(1000000, 1500, fname)
	} else {

		fmt.Println("Setup plumbing...")

		lineChan := make(chan string)
		parsedChans := make([]chan ParsedData, NO_OF_AGGREGATOR_WORKERS)
		resultChan := make(chan AggregatorResult, NO_OF_AGGREGATOR_WORKERS)

		// Create aggregator channels
		for i := range parsedChans {
			parsedChans[i] = make(chan ParsedData, 100)
		}

		// Start aggregators
		var wgAggregators sync.WaitGroup
		for i := 0; i < NO_OF_AGGREGATOR_WORKERS; i++ {
			wgAggregators.Add(1)
			go aggregatorWorker(i, parsedChans[i], resultChan, &wgAggregators)
		}

		// Start parsers
		var wgParsers sync.WaitGroup
		for i := 0; i < NO_OF_PARSER_WORKERS; i++ {
			wgParsers.Add(1)
			go parserWorker(i, lineChan, parsedChans, NO_OF_AGGREGATOR_WORKERS, &wgParsers)
		}

		fmt.Println("Starting pipeline...")
		startTime := time.Now()

		// Reader
		err := util.ReadFileLines(fname, lineChan)
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
		finalMap := make(map[string]float64)
		var totalStats AggregatorStats
		for res := range resultChan {
			fmt.Printf("Aggregator %d stats: %d lines, %d keys, total %.2f\n",
				res.ID, res.Stats.LinesProcessed, res.Stats.UniqueKeys, res.Stats.TotalValue)

			for k, v := range res.Data {
				finalMap[k] += v
			}

			totalStats.LinesProcessed += res.Stats.LinesProcessed
			totalStats.UniqueKeys += res.Stats.UniqueKeys // rough count, may have overlap
			totalStats.TotalValue += res.Stats.TotalValue
		}

		// Output
		fmt.Println("\n Final aggregated results:")
		for k, v := range finalMap {
			fmt.Printf("%s: %.2f\n", k, v)
		}

		elapsed := time.Since(startTime)
		fmt.Printf("\nDone in %s. Processed %d lines, approx. %d unique keys, total sum %.2f\n",
			elapsed, totalStats.LinesProcessed, len(finalMap), totalStats.TotalValue)
	}
}

func parserWorker(id int, lines <-chan string, parsedChans []chan ParsedData, shardCount int, wg *sync.WaitGroup) {
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

		data := ParsedData{Key: key, Value: value}
		shard := hashKey(key) % shardCount
		parsedChans[shard] <- data
	}
}

func aggregatorWorker(id int, input <-chan ParsedData, out chan<- AggregatorResult, wg *sync.WaitGroup) {
	defer wg.Done()

	localMap := make(map[string]float64)
	var stats AggregatorStats

	for data := range input {
		localMap[data.Key] += data.Value
		stats.LinesProcessed++
		stats.TotalValue += data.Value
	}

	stats.UniqueKeys = len(localMap)

	out <- AggregatorResult{
		ID:    id,
		Data:  localMap,
		Stats: stats,
	}
}

func hashKey(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}

/*
hashmap := make(map[string]int)
hashmap["A"] = 25
value, exists := hashmap["A"]
isEmpty := len(hashmap) == 0
for key, value := range hashmap {
        fmt.Printf("%s -> %d\n", key, value)
}
toSlice := make([]int, 0, len(s.data))
    for key := range s.data {
        result = append(result, key)
}
delete(hashmap, "A")
*/
