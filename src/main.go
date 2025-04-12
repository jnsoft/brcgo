package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/brcgo/src/models"
	"github.com/brcgo/src/util"
)

const (
	NO_OF_PARSER_WORKERS     = 4
	NO_OF_AGGREGATOR_WORKERS = 4
	ERROR                    = "❌ Error reading file"
	WARNING                  = "⚠️ Warning"
	DONE                     = "✅ Done"
)

func main() {
	fname := "testfile_10000000.tmp"

	if false {
		util.GenerateFile(1000000, 1500, fname)
		return
	}

	Naive(fname, false)
	//pipelines.WorkerpoolPipeline(fname, 10, false)
	//pipelines.ReadParseAggregatePipeline(fname, NO_OF_PARSER_WORKERS, NO_OF_AGGREGATOR_WORKERS, false)

}

func Naive(fname string, verbose bool) error {

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	startTime := time.Now()

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	resultMap := make(map[string]models.StationData)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := models.ParseLine(scanner.Text())
		aggregated, exists := resultMap[data.Key]
		if !exists {
			resultMap[data.Key] = models.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			resultMap[data.Key] = models.StationData{
				Min:   math.Min(data.Value, aggregated.Min),
				Max:   math.Max(data.Value, aggregated.Max),
				Sum:   data.Value + aggregated.Sum,
				Count: aggregated.Count + 1,
			}
		}

	}

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

	return scanner.Err()
}
