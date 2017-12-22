
Benchmark 测试结果如下：

```shell
➜  bigcache git:(master) ✗ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gu/libs/store/bigcache
BenchmarkCacheDB_Get-4                    500000              2155 ns/op
BenchmarkCacheServiceProvider_Get-4      2000000               767 ns/op
PASS
ok      github.com/fengyfei/gu/libs/store/bigcache      3.486s
```