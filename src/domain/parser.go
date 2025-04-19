package domain

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

func ParseStringFloat(s string) (StringFloat, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid line: %s", s))
	}
	key := strings.TrimSpace(parts[0])
	valStr := strings.TrimSpace(parts[1])

	value, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		panic(fmt.Sprintf(" Failed to parse float: %s", s))
	}

	return StringFloat{Key: key, Value: value}, nil
}

func ParseStringInt(s string) StringInt {
	lg := len(s)
	var key string
	var ix int
	for ix = 0; ix < lg; ix++ {
		if s[ix] == ';' {
			break
		}
	}
	neg := 0
	if s[ix+1] == '-' {
		neg = 1
	}
	fac := 1
	n := 0
	for i := lg - 1; i > ix+neg; i-- {
		if s[i] != '.' {
			n += fac * int(s[i])
			fac *= 10
		}
	}

	key = s[:ix]
	return StringInt{Key: key, Value: n}
}

func ParseBytesInt(bs []byte) BytesInt {
	lg := len(bs)

	var ix int
	// find split
	for ix = 0; ix < lg; ix++ {
		if bs[ix] == ASCII_SEMICOLON {
			break
		}
	}
	neg := 0
	if bs[ix+1] == ASCII_MINUS {
		neg = 1
	}
	// read number
	fac := 1
	n := 0
	for i := lg - 1; i > ix+neg; i-- {
		if bs[i] != ASCII_DOT {
			n += fac * int(ASCII_ZERO-bs[i])
			fac *= 10
		}
	}

	key := bs[:ix]
	return BytesInt{Key: key, Value: n}
}
