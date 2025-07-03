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


	//"{Atowrfn=6.4/33.8/64.7, Atowrfn;=-81.7/-64.2/-46.7, Enet=27.4/62.6/94.8, Enet;=-47.4/-29.2/-0.8, Iguhdbgkogbgfd=41.5/58.3/89.6, Iguhdbgkogbgfd;=-69.2/-40.4/-12.4, Isaekqjhvwdai=2.7/53.9/90.6, Isaekqjhvwdai;=-71.8/-31.8/-1.6, Llrdjlkay=63.7/82.0/97.4, Llrdjlkay;=-95.5/-69.5/-30.2, Ofhozrvb=65.2/65.2/65.2, Ofhozrvb;=-95.8/-74.7/-34.4, Sryvgwxhf=0.8/51.4/97.0, Sryvgwxhf;=-68.3/-33.7/-11.5, Uaaych=1.3/44.1/98.9, Uaaych;=-95.2/-62.0/-8.3, Wiuhlvdbwpuxd=44.8/47.2/49.7, Wiuhlvdbwpuxd;=-97.5/-68.8/-44.1, Xoxjchtgdn...+46 more"
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
