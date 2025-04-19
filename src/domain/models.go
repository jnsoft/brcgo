package domain

import (
	"fmt"
)

type StringFloat struct {
	Key   string
	Value float64
}

type StringInt struct {
	Key   string
	Value int
}

type BytesInt struct {
	Key   []byte
	Value int
}

type StationData struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

type StationDataInt struct {
	Min   int
	Max   int
	Sum   int
	Count int
}

func (s StationData) String() string {
	return fmt.Sprintf("%.2f/%.2f/%.2f", s.Min, s.Sum/float64(s.Count), s.Max)
}

func (s StationDataInt) String() string {
	return fmt.Sprintf("%.2f/%.2f/%.2f", float64(s.Min)/10, float64(s.Sum)/float64(s.Count*10), float64(s.Max)/10)
}
