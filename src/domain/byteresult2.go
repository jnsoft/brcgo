package domain

import (
	"fmt"
	"sort"
	"strings"
)

type ByteResult2 struct {
	stations map[int]*ByteStation
	inputs   int
}

func NewByteResult2() *ByteResult2 {
	return &ByteResult2{
		stations: make(map[int]*ByteStation),
		inputs:   0,
	}
}

func (r *ByteResult2) NoOfStations() int {
	return len(r.stations)
}

func (r *ByteResult2) NoOfInputs() int {
	return r.inputs
}

func (r *ByteResult2) GetStations() int {
	return len(r.stations)
}

func (r *ByteResult2) Add(reading ByteStationReading) {
	key := reading.HashCodeSimple()
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

func (r *ByteResult2) Merge(r2 *ByteResult2) {
	for key, station := range r2.stations {
		r.inputs += r2.inputs
		existingStation, exists := r.stations[key]
		if !exists {
			r.stations[key] = &ByteStation{
				StationId: station.StationId,
				Min:       station.Min,
				Max:       station.Max,
				Sum:       station.Sum,
				Count:     station.Count,
			}
		} else {
			existingStation.Sum += station.Sum
			if existingStation.Min > station.Min {
				existingStation.Min = station.Min
			}
			if existingStation.Max < station.Max {
				existingStation.Max = station.Max
			}
			existingStation.Count += station.Count
		}
	}
}

func (r *ByteResult2) String() string {
	return fmt.Sprintf("Processed %d stations.", len(r.stations))
}

func (r *ByteResult2) GetSortedResults() string {
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

func (r *ByteResult2) PrintResults() {
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
