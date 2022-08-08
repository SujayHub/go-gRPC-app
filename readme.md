###### run the following command to resolve module dependencies
```shell
protoc --proto_path=. \
--go_out=. \
--go_opt=paths=import \
--go_opt=module=github.com/sujayhub/go-gRPC-app \
--go-grpc_out=. \
--go-grpc_opt=paths=import \
--go-grpc_opt=module=github.com/sujayhub/go-gRPC-app \
protos/*.proto
```