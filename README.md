# color-it

Implementation for the ["color-it"](https://www.sfeir.com/fr/battle-algo/) contest.

## Building

### Build the binary using the Go SDK (preferred)

```bash
# Install the application dependencies.
go mod download

# Build the application.
go build

# Test it.
./color-it samples/30_30_3-1.csv
```

### Build a Docker image

```bash
# Build the Docker image
docker build -t color-it .

# Test it.
docker run --rm -v $(pwd)/samples/30_30_3-1.csv:/data/input.csv color-it /data/input.csv
```

## Usage

The only required parameter to run the program is the input CSV file to process; it is passed as a positional argument.
Some other optional arguments can be provided to control the program behavior but the default values should be used for
the contest.

```bash
Usage of ./color-it:
  -check-square
        Check whether the board is a square after loading it (default true)
  -debug
        Enable the debug logs
  -impl string
        Name of the algorithm implementation to execute (default "deep-search")
  -timeout int
        Timeout in seconds of the execution (default 115)
```

### Output

The best solution found is printed on stdout, one step per line at the end of the program execution, for example:
```bash
1
2
0
2
```

The file `out.csv` is also created at the end of the program execution and contains the solution (same format).

## Results

| Sample        | Deep search                                                                            |
|---------------|----------------------------------------------------------------------------------------|
| 12_12_4-1.csv | **0m0,255s** <br />nb-steps=12 <br />solution=[1,0,3,1,3,0,1,0,2,3,1,0]                |
| 12_12_5-1.csv | **0m0,391s** <br />nb-steps=14 <br />solution=[0,4,3,4,1,2,3,4,2,0,3,4,2,1]            |
| 12_12_6-1.csv | **1m29,299s** <br />nb-steps=19 <br />solution=[0,2,4,5,3,1,0,5,2,3,0,5,4,2,3,5,4,0,1] |
| 15_15_3-1.csv | **0m0,042s** <br />nb-steps=9 <br />solution=[1,2,0,1,0,2,0,1,2]                       |
| 15_15_3-2.csv | **0m0,045s** <br />nb-steps=9 <br />solution=[2,0,1,0,2,0,1,2,0]                       |
| 15_15_4-1.csv | **0m4,270s** <br />nb-steps=16 <br />solution=[3,2,0,1,3,1,0,1,3,0,2,3,1,2,0,3]        |
| 15_15_5-1.csv | **1m19,528s** <br />nb-steps=17 <br />solution=[3,0,2,3,4,1,2,1,3,0,2,4,0,1,3,2,4]     |
| 15_15_6-1.csv | **timeout** <br />nb-steps=24 <br />solution=[...]                                     |
| 20_20_3-1.csv | **0m0,049ss**  <br />nb-steps=11 <br />solution=[0,2,1,2,0,1,2,1,0,2,1]                |
| 20_20_4-1.csv | **1m27,674s** <br />nb-steps=19 <br />solution=[3,1,0,1,3,2,3,0,2,1,0,3,1,2,3,1,0,2,3] |
| 20_20_5-1.csv | **timeout** <br />nb-steps=26 <br />solution=[...]                                     |
| 20_20_6-1.csv | **timeout** <br />nb-steps=34 <br />solution=[...]                                     |
| 30_30_3-1.csv | **0m4,027s** <br />nb-steps=19 <br />solution=[1,2,0,2,1,2,0,1,2,0,1,2,1,2,1,0,2,1,0]  |
| 30_30_4-1.csv | **timeout** <br />nb-steps=28 <br />solution=[...]                                     |
| 30_30_5-1.csv | **timeout** <br />nb-steps=38 <br />solution=                                          |
| 30_30_6-1.csv | **timeout** <br />nb-steps=51 <br />solution=[...]                                     |

## Profiling

Use go test benchmark feature to generate the profiling files:

```bash
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
```

Use the [pprof](https://github.com/google/pprof) tool to visualize the profiling results with pprof:

```bash
go tool pprof -http=":" cpu.prof
```

## TODO

Parallelize the deep search implementation using multiple Goroutines.

Try new implementation:

1. start from the initial solution (of length L)
1. start from level L-2 and test the other colors to see if there is a better solution at this level
1. go upward until reaching the tree root
