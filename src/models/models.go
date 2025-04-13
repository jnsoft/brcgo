package models

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ASCII_SEMICOLON = 59 // ';'
	ASCII_MINUS     = 45 // '-'
	ASCII_DOT       = 46 // '.'
	ASCII_ZERO      = 48 // 0=48, 1=49...
)

type ParsedData struct {
	Key   string
	Value float64
}

type ParsedDataInt struct {
	Key   string
	Value int
}

type ParsedDataByteInt struct {
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
	for i := lg - 1; i > ix+neg; i-- {
		if line[i] != '.' {
			n += fac * int(line[i])
			fac *= 10
		}
	}

	key = line[:ix]
	return ParsedDataInt{Key: key, Value: n}
}

func ParseByteLineInt(line []byte) ParsedDataInt {
	lg := len(line)

	var ix int
	// find split
	for ix = 0; ix < lg; ix++ {
		if line[ix] == ASCII_SEMICOLON {
			break
		}
	}
	neg := 0
	if line[ix+1] == ASCII_MINUS {
		neg = 1
	}
	// read number
	fac := 1
	n := 0
	for i := lg - 1; i > ix+neg; i-- {
		if line[i] != ASCII_DOT {
			n += fac * int(ASCII_ZERO-line[i])
			fac *= 10
		}
	}

	key := string(line[:ix])
	return ParsedDataInt{Key: key, Value: n}
}
