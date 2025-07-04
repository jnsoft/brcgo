package domain

func Aggregate(data StringFloat, hashmap *map[string]*StationData) {
	aggregated, exists := (*hashmap)[data.Key]
	if !exists {
		(*hashmap)[data.Key] = &StationData{
			Min:   data.Value,
			Max:   data.Value,
			Sum:   data.Value,
			Count: 1,
		}
	} else {
		if data.Value < aggregated.Min {
			aggregated.Min = data.Value
		} else if data.Value > aggregated.Max {
			aggregated.Max = data.Value
		}
		aggregated.Sum += data.Value
		aggregated.Count++
	}
}
