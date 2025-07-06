package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/brcgo/src/domain"
	"github.com/brcgo/src/pipelines"
	"github.com/brcgo/src/util"
	"github.com/brcgo/src/workers"
	"github.com/jnsoft/jngo/misc"
	"github.com/jnsoft/jngo/pipeline"
	"github.com/jnsoft/jngo/profiling"
)

const (
	NO_OF_PARSER_WORKERS     = 4
	NO_OF_AGGREGATOR_WORKERS = 4
	MAX_NO_OF_ROWS           = 1000000000
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
	fname := flag.String("f", "", "The name of the file to read")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	no_of_pallell := flag.Int("p", 1, "Maximum number of concurrent threads")
	generate := flag.Bool("g", false, "Create test file")
	no_of_rows := flag.Int("r", 100, "Number of rows to generate")
	no_of_stations := flag.Int("s", 10, "Number of stations in generated file")

	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix(time.Now().Format(time.RFC3339) + " ")

	if *fname == "" {
		log.Fatal("Filename is required: -f <file_name>")
	} else if *generate {
		if *no_of_rows <= 0 {
			log.Fatal("Number of rows must be greater than 0")
		}
		if *no_of_stations <= 0 {
			log.Fatal("Number of stations must be greater than 0")
		}
		if *no_of_rows > MAX_NO_OF_ROWS {
			*no_of_rows = MAX_NO_OF_ROWS
		}
		log.Printf("Generating file with %d rows and %d stations\n", *no_of_rows, *no_of_stations)
		util.GenerateFile(*no_of_rows, *no_of_stations, *fname)
		log.Printf("File generated: %s\n", *fname)
		return

	} else {
		if _, err := os.Stat(*fname); os.IsNotExist(err) {
			log.Fatalf("File does not exist: %s", *fname)
		}

	}

	if verbose != nil && *verbose {
		log.Println("Verbose mode enabled")
	}
	log.Printf("Using file %s", *fname)
	log.Printf("Using %d parallel workers", *no_of_pallell)

	//pipelines.Naive(*fname)

	//pipelines.NaiveBytes(*fname, *no_of_pallell)

	pipelines.ParallellMmap(*fname, *no_of_pallell)

	pipelines.SimpleMmap(*fname)

	//TestChannel2()

	//TestContext()

	// RunPipeline(fname, verbose)
	//RunPipeline2(fname, verbose)

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

func TestContext() {
	timeOutContext, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	req, err := http.NewRequestWithContext(timeOutContext, http.MethodGet, "https://www.google.com", nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

// simulate work
func DoWork(i int) int {
	time.Sleep(time.Millisecond * 500)
	return i
}

func WaitSeconds(seconds int) bool {
	time.Sleep(time.Second * time.Duration(seconds))
	return true
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

	profiling.ProfileFunction("Pipelinebuilder", PROF_FNAME, func() (interface{}, error) {
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

	profiling.ProfileFunction("Ideomatic pipeline", PROF_FNAME, func() (interface{}, error) {
		pipelines.IdeomotaticPipeline[domain.StringFloat, domain.StationData](fname,
			domain.ParseStringFloat,
			collector,
			printer,
			verbose,
		)
		return 0, nil
	})
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
