#### run all the commands from checked-out project directory
###### fetch all the dependencies by running the following command
```shell
go mod tidy
```
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

###### run the following command to start the server
```shell
go run server/main.go
```

###### run the following command to start the client
```shell
go run client/main.go
```
