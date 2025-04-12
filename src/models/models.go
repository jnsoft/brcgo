package models

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (s StationData) String() string {
	return fmt.Sprintf("%.2f/%.2f/%.2f", s.Min, s.Sum/float64(s.Count), s.Max)
}

func ParseLine(line string) ParsedData {
	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid line: %s", line))
	}
	key := strings.TrimSpace(parts[0])
	valStr := strings.TrimSpace(parts[1])

	value, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		panic(fmt.Sprintf(" Failed to parse float: %s", line))
	}

	return ParsedData{Key: key, Value: value}
}
