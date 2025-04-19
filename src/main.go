package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/brcgo/src/domain"
	"github.com/brcgo/src/misc"
	"github.com/brcgo/src/pipeline"
	"github.com/brcgo/src/pipelines"
	"github.com/brcgo/src/util"
	"github.com/brcgo/src/workers"
)

const (
	NO_OF_PARSER_WORKERS     = 4
	NO_OF_AGGREGATOR_WORKERS = 4
	ERROR                    = "❌ Error reading file"
	WARNING                  = "⚠️ Warning"
	DONE                     = "✅ Done"
	PROF_FNAME               = "cpu_profile.prof"
)

func main() {
	fname := "testfile_1000000.tmp"
	verbose := false

	if false {
		util.GenerateFile(1000000, 1500, fname)
		return
	}

	hashmap := make(map[string]*domain.StationData)

	pipeline.Pipeline[string](
		func(out chan<- string) error {
			return workers.GetLines(fname, out)
		},
		[]pipeline.Stage[string]{
			func(in <-chan string) <-chan domain.StringFloat {
				return pipelines.ParallelMapStage[string, domain.StringFloat](8, domain.ParseStringFloat)(in)
			},
			func(in <-chan domain.StringFloat) <-chan domain.StringFloat {
				return pipelines.ParallelCollectStage(4, func(d domain.StringFloat) {
					domain.Aggregate(d, &hashmap)
				})(in)
			},
		},
		func() {
			domain.PrintResult(&hashmap, verbose)
		},
	)

	collector := func(data domain.StringFloat) {
		domain.Aggregate(data, &hashmap)
	}
	printer := func() {
		domain.PrintResult(&hashmap, verbose)
	}

	misc.ProfileFunction("Ideomatic pipeline", PROF_FNAME, func() (interface{}, error) {
		pipelines.IdeomotaticPipeline[domain.StringFloat, domain.StationData](fname,
			domain.ParseStringFloat,
			collector,
			printer,
			verbose,
		)
		return 0, nil
	})

	Naive(fname, verbose)

	//misc.ProfileFunction("Naive int", PROF_FNAME, func() (interface{}, error) {
	//	return NaiveInt(fname, false), nil
	//})

	misc.ProfileFunction("Naive int 2", PROF_FNAME, func() (interface{}, error) {
		return NaiveInt2(fname, false), nil
	})

	//pipelines.WorkerpoolPipeline(fname, 10, false)
	//pipelines.ReadParseAggregatePipeline(fname, NO_OF_PARSER_WORKERS, NO_OF_AGGREGATOR_WORKERS, false)

}

func Naive(fname string, verbose bool) error {
	startTime := time.Now()

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	resultMap := make(map[string]domain.StationData)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data, _ := domain.ParseStringFloat(scanner.Text())
		aggregated, exists := resultMap[data.Key]
		if !exists {
			resultMap[data.Key] = domain.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			resultMap[data.Key] = domain.StationData{
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

func NaiveInt2(fname string, verbose bool) error {

	startTime := time.Now()
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	resultMap := make(map[string]domain.StationDataInt)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := domain.ParseStringInt(scanner.Text())
		aggregated, exists := resultMap[data.Key]
		if !exists {
			resultMap[data.Key] = domain.StationDataInt{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			resultMap[data.Key] = domain.StationDataInt{
				Min:   misc.Min(data.Value, aggregated.Min),
				Max:   misc.Max(data.Value, aggregated.Max),
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
