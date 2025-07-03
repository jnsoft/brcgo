package domain

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type ByteResult struct {
	stations map[int]*ByteStation
	inputs   int
	mu       sync.Mutex
}

func NewByteResult() *ByteResult {
	return &ByteResult{
		stations: make(map[int]*ByteStation),
		inputs:   0,
	}
}

func (r *ByteResult) NoOfStations() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.stations)
}

func (r *ByteResult) NoOfInputs() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.inputs
}

func (r *ByteResult) GetStations() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.stations)
}

func (r *ByteResult) Add(reading ByteStationReading) {
	key := reading.HashCodeSimple()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inputs++
	station, exists := r.stations[key]
	if !exists {
		r.stations[key] = &ByteStation{
			StationId: reading.StationId,
			Min:       reading.Temperature,
			Max:       reading.Temperature,
			Sum:       int64(reading.Temperature),
			Count:     1,
		}
	} else {
		station.Sum += int64(reading.Temperature)
		if station.Min > reading.Temperature {
			station.Min = reading.Temperature
		}
		if station.Max < reading.Temperature {
			station.Max = reading.Temperature
		}
		station.Count++
	}
}

func (r *ByteResult) String() string {
	return fmt.Sprintf("Processed %d stations.", len(r.stations))
}

func (r *ByteResult) GetSortedResults() string {
	stations := make([]*ByteStation, 0, len(r.stations))
	for _, v := range r.stations {
		stations = append(stations, v)
	}
	sort.Slice(stations, func(i, j int) bool {
		return stations[i].StationName() < stations[j].StationName()
	})

	var sb strings.Builder
	sb.WriteByte('{')
	for i, s := range stations {
		sb.WriteString(s.String())
		if i < len(stations)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

func (r *ByteResult) PrintResults() {
	fmt.Println("StationId\tMin\tMax\tSum\tCount")
	fmt.Println("---------------------------------------")
	// Sort by StationName
	stations := make([]*ByteStation, 0, len(r.stations))
	for _, v := range r.stations {
		stations = append(stations, v)
	}
	sort.Slice(stations, func(i, j int) bool {
		return stations[i].StationName() < stations[j].StationName()
	})
	for _, s := range stations {
		fmt.Println(s)
	}
}
