package pipelines

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/brcgo/src/domain"
)

func Naive(fname string) error {
	startTime := time.Now()

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	resultMap := make(map[string]domain.StationData)
	cnt := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cnt++
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

	//if verbose {
	//	fmt.Println("\n Final aggregated results:")
	//	for _, k := range keys {
	//		fmt.Printf("%s=%s\n", k, resultMap[k].String())
	//	}
	//}

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone in %s. Processed %d lines, %d unique keys\n",
		elapsed, cnt, len(resultMap))

	return scanner.Err()
}
