```
go test -bench=BenchmarkCoins -run=^a

go test -bench=BenchmarkCoins -run=^a -benchtime=10s

goos: linux
goarch: amd64
pkg: github.com/Lunkov/go-ecos-client
cpu: 12th Gen Intel(R) Core(TM) i7-12700H
BenchmarkCoins/PROCESS(1)_-20         	   44626	    266507 ns/op
BenchmarkCoins/PROCESS(2)_-20         	   49509	    245843 ns/op
BenchmarkCoins/PROCESS(4)_-20         	   48970	    244578 ns/op
BenchmarkCoins/PROCESS(8)_-20         	   50565	    243124 ns/op
count transactions        = 233964
work time                 = 58.523822 seconds
Work time per transaction = 0.000250 seconds
Make transactions per second = 3997.756664 seconds



```
