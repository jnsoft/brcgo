package domain

import (
	"fmt"
	"sort"
)

func PrintResult(hashmap *map[string]*StationData, verbose bool) {
	keys := make([]string, 0, len(*hashmap))
	for k := range *hashmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if verbose {
		fmt.Println("\n Final aggregated results:")
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, (*hashmap)[k].String())
		}
	}
	fmt.Printf("\n%d unique keys\n",
		len(*hashmap))
}
