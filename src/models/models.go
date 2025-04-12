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

type ParsedDataInt struct {
	Key   string
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

func ParseLineGeneric(line string) ParsedData {
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

func ParseLine(line string) ParsedData {
	lg := len(line)
	var key string
	var ix int
	for ix = 0; ix < lg; ix++ {
		if line[ix] == ';' {
			break
		}
	}
	key = line[:ix]
	val := line[ix+1:]

	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf(" Failed to parse float: %s", line))
	}

	return ParsedData{Key: key, Value: value}
}

func ParseLineInt(line string) ParsedDataInt {
	lg := len(line)
	var key string
	var ix int
	for ix = 0; ix < lg; ix++ {
		if line[ix] == ';' {
			break
		}
	}
	neg := 0
	if line[ix+1] == '-' {
		neg = 1
	}
	fac := 1
	n := 0
	for i := lg - 1; i > ix+neg; ix++ {
		if line[i] != '.' {
			n += fac * line[i]
		}
	}

	key = line[:ix]
	val := line[ix+1:]

	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf(" Failed to parse float: %s", line))
	}

	return ParsedDataInt{Key: key, Value: value}
}
