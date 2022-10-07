# color-it

Implementation for the ["color-it"](https://www.sfeir.com/fr/battle-algo/) contest.

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
1. start from level L-2 and check if the other colors to see if there is a better solution
1. go upward until reaching the tree root
