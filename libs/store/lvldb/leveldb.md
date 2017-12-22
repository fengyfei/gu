
Benchmark 测试结果如下：

```shell
➜  lvldb git:(master) ✗ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gu/libs/store/lvldb
BenchmarkDatabase_Put-4           100000             17219 ns/op
BenchmarkDatabase_Get-4            10000            159853 ns/op
PASS
ok      github.com/fengyfei/gu/libs/store/lvldb 3.885s
```