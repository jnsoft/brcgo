package domain

import "fmt"

type ByteStation struct {
	StationId   []byte
	Sum         int64
	Min         int
	Max         int
	Count       int
	stationName string
}

func NewByteStation() *ByteStation {
	return &ByteStation{
		StationId:   []byte{},
		Sum:         0,
		Count:       0,
		stationName: "",
	}
}

func (b *ByteStation) StationName() string {
	if b.stationName == "" {
		if len(b.StationId) == 0 {
			b.stationName = ""
		} else {
			b.stationName = string(b.StationId) // cache the station name
		}
	}
	return b.stationName
}

func (b *ByteStation) String() string {
	return fmt.Sprintf("%s=%.1f/%.1f/%.1f",
		b.StationName(),
		float64(b.Min)/10,
		b.averageTemperature(),
		float64(b.Max)/10,
	)
}

func (b *ByteStation) averageTemperature() float64 {
	if b.Count > 0 {
		return float64(b.Sum) / float64(b.Count*10)
	}
	return 0.0
}
