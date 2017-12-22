A 通过 vice 向 B 发送信息，benchmark 多次测试结果：

1.

```shell
count --->>> 24158
  100000	     11779 ns/op
PASS
ok  	github.com/fengyfei/gopher/vice/chat/clienti	1.622s
```

2.

```shell
count --->>> 16956
  200000	      5582 ns/op
PASS
ok  	github.com/fengyfei/gopher/vice/chat/clienti	1.508s
```

3.

```shell
count --->>> 27205
  300000	      3950 ns/op
PASS
ok  	github.com/fengyfei/gopher/vice/chat/clienti	1.473s
```

4.

```shell
count --->>> 22381
  100000	     11327 ns/op
PASS
ok  	github.com/fengyfei/gopher/vice/chat/clienti	1.492s
```

