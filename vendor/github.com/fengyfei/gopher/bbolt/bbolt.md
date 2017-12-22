
Benchmark 测试结果如下：

```shell
➜  bbolt git:(master) ✗ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gopher/bbolt
BenchmarkUserServiceProvider_Put-4          3000            375388 ns/op
BenchmarkUserServiceProvider_Get-4       1000000              1899 ns/op
PASS
ok      github.com/fengyfei/gopher/bbolt        3.120s
```