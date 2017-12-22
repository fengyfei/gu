
Benchmark 测试结果如下：

```shell
➜  badger git:(master) ✗ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/fengyfei/gu/libs/store/badger
BenchmarkBadgerDB_Set-4            10000            211903 ns/op
BenchmarkBadgerDB_Get-4           500000              2401 ns/op
PASS
ok      github.com/fengyfei/gu/libs/store/badger        3.528s
```