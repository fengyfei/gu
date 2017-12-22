
Benchmark 测试结果如下：

```shell
➜  goleveldb git:(master) ✗ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gopher/goleveldb
BenchmarkAccountServiceProvider_Create-4             100          18458620 ns/op
BenchmarkAccountServiceProvider_GetRandom-4          100          13149320 ns/op
PASS
ok      github.com/fengyfei/gopher/goleveldb    3.209s
```