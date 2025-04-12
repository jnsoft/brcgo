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

## Start
```
go mod init github.com/goEuler
```

## Run and Test
```
go build -v ./...
go test -v ./...
```


```
