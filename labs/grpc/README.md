# gRPC

## References

- [gRPC](https://grpc.io/docs/tutorials/basic/go.html)
- [Protobuf 3](https://developers.google.com/protocol-buffers/docs/proto3)
- [Go gRPC Middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)

## protoc

```sh
cd sample
protoc -I greeter/ greeter/greeter.proto --go_out=plugins=grpc:greeter
```
