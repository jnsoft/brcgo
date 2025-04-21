package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
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

var (
	hashmap = make(map[string]*domain.StationData)
	mu      sync.Mutex
)

func main() {
	fname := "testfile_1000000.tmp"
	verbose := false
	if verbose {
		fmt.Println(verbose)
	}

	if false {
		util.GenerateFile(1000000, 1500, fname)
		return
	}

	TestChannel2()

	// RunPipeline(fname, verbose)
	//RunPipeline2(fname, verbose)
	//Naive(fname, verbose)

	//misc.ProfileFunction("Naive int", PROF_FNAME, func() (interface{}, error) {
	//	return NaiveInt(fname, false), nil
	//})

	//misc.ProfileFunction("Naive int 2", PROF_FNAME, func() (interface{}, error) {
	//	return NaiveInt2(fname, false), nil
	//})

	//pipelines.WorkerpoolPipeline(fname, 10, false)
	//pipelines.ReadParseAggregatePipeline(fname, NO_OF_PARSER_WORKERS, NO_OF_AGGREGATOR_WORKERS, false)

}

func WaitGroupExample() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i) // falls func with value i
	}

	wg.Wait() // blocks until all goroutines finish
}

func MutexExample() {
	var mu sync.Mutex
	count := 0

	for i := 0; i < 5; i++ {
		go func() {
			mu.Lock()
			count++ // safe from race condition
			mu.Unlock()
		}()
	}
}

func TestChannel1() {
	dataChan := make(chan int)

	go func() {
		for i := range 100 {
			dataChan <- i
		}
		close(dataChan)
	}()

	for n := range dataChan {
		fmt.Printf("n=%d\n", n)
	}
}

func TestChannel2() {

	dataChan := make(chan int)

	go func() {
		wg := sync.WaitGroup{}

		for i := range 100 {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				res := DoWork(j)
				dataChan <- res
			}(i)
		}
		wg.Wait()
		close(dataChan)
	}()

	for n := range dataChan {
		fmt.Printf("n=%d\n", n)
	}
}

// simulate work
func DoWork(i int) int {
	time.Sleep(time.Millisecond * 500)
	return i
}

func RunPipeline(fname string, verbose bool) {
	start := time.Now()

	pb := pipeline.FromSource(func(out chan<- string) error {
		return workers.GetLines(fname, out)
	})

	pb2 := pipeline.Then(pb, pipeline.ParallelMapStage[string, domain.StringFloat](8, domain.ParseStringFloat))

	pb3 := pipeline.Then(pb2, pipeline.ParallelDoStage[domain.StringFloat](8, func(data domain.StringFloat) {
		mu.Lock()
		defer mu.Unlock()

		aggregated, exists := hashmap[data.Key]
		if !exists {
			hashmap[data.Key] = &domain.StationData{
				Min:   data.Value,
				Max:   data.Value,
				Sum:   data.Value,
				Count: 1,
			}
		} else {
			if data.Value < aggregated.Min {
				aggregated.Min = data.Value
			} else if data.Value > aggregated.Max {
				aggregated.Max = data.Value
			}
			aggregated.Sum += data.Value
			aggregated.Count++
		}
	}))

	misc.ProfileFunction("Pipelinebuilder", PROF_FNAME, func() (interface{}, error) {
		pipeline.Run(pb3, func() {
			domain.PrintResult(&hashmap, verbose)
			fmt.Printf("Processed %d keys in %v\n", len(hashmap), time.Since(start))
		})
		return len(hashmap), nil
	})
}

func RunPipeline2(fname string, verbose bool) {
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
