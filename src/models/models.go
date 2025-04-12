package models

import "fmt"

type ParsedData struct {
	Key   string
	Value float64
}

type StationData struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func (p StationData) String() string {
	return fmt.Sprintf("%.2f/%.2f/%.2f", p.Min, p.Sum/float64(p.Count), p.Max)
}
