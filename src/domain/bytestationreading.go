package domain

type ByteStationReading struct {
	StationId   []byte
	Temperature int
}

func NewByteStationReading() ByteStationReading {
	return ByteStationReading{
		StationId:   []byte{},
		Temperature: 0,
	}
}

func NewByteStationReadingFromBytes(bs []byte) ByteStationReading {
	lg := len(bs)
	var ix int

	// find split
	for ix = 0; ix < lg; ix++ {
		if bs[ix] == ASCII_SEMICOLON {
			break
		}
	}

	// check sign
	neg := 0
	if bs[ix+1] == ASCII_MINUS {
		neg = 1
		ix++
	}

	// read number
	fac := 1
	n := 0
	for i := lg - 1; i > ix; i-- {
		if bs[i] != ASCII_DOT {
			n += fac * int(bs[i]-ASCII_ZERO)
			fac *= 10
		}
	}
	if neg == 1 {
		n = -n
	}

	return ByteStationReading{
		StationId:   bs[:ix],
		Temperature: n,
	}
}

func (r ByteStationReading) HashCode() int {
	const prime = 16777619
	hash := uint32(2166136261)
	for _, b := range r.StationId {
		hash = (hash ^ uint32(b)) * prime
	}
	return int(hash)
}

func (r ByteStationReading) HashCodeSimple() int {
	hash := 17
	for _, b := range r.StationId {
		hash = hash*31 + int(b)
	}
	return hash
}
