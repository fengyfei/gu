Benchmark 测试结果如下：

```shell
bigcache git:(master) ✗ go test -bench=.  
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gopher/bigcache
BenchmarkCacheServiceProvider_SetOne-4            500000              2384 ns/op
BenchmarkCacheServiceProvider_GetOne-4           1000000              1016 ns/op
PASS
ok      github.com/fengyfei/gopher/bigcache     2.328s
```