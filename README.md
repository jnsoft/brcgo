# Billion Rows Challenge in Go

## Challenge
The challenge is to read a file with 1 billion lines, aggregating the information contained in each line and print a report with the result. Each line contains a weather station name and a temperature reading in the format `<station name>;<temperature>`, Station name may have spaces and other special characters excluding ;, and the temperature is a floating-point number ranging from -99.9 to 99.9 with precision limited to one decimal point. The expected output format is `{<station name>=<min>/<mean/<max>, ...}`  

Example of a measurement file:
```
Yellowknife;16.0
Entebbe;32.9
Porto;24.4
Vilnius;12.4
Fresno;7.9
Maun;17.5
Panama City;39.5
...
```

Example of expected output:
```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, ...}
```

## Attemp 1  
* Reader: sends lines to lineChan
* Parsers: convert lines to (key, float) and send to correct aggregator (based on hash of key)
* Aggregators (each owns a set of keys):
 * Track local stats
 * Send final map + stats to main thread
* Main thread: Merges all maps, prints stats and results



## Start
```
read -s username
git config --global user.email $username@users.noreply.github.com && git config --global user.name $username

go mod init github.com/jnsoft/brcgo

GOPRIVATE=github.com/jnsoft/jngo go get github.com/jnsoft/jngo@latest
```

## Run and Test
```
go build -v ./...
go test -v ./...
cd src/cache
go test -bench BenchmarkSimpleCache -count=2

go tool pprof -http 127.0.0.1:8080 ./cpu_profile.prof
```


```


### Extra

```
hashmap := make(map[string]int)
hashmap["A"] = 25
value, exists := hashmap["A"]
isEmpty := len(hashmap) == 0
for key, value := range hashmap {
        fmt.Printf("%s -> %d\n", key, value)
}
toSlice := make([]int, 0, len(s.data))
    for key := range s.data {
        result = append(result, key)
}
delete(hashmap, "A")
```
