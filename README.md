# grpc-gateway-example

grpc代理简单示例，http协议转grpc协议。

client ---http1.1---> proxy ---grpc---> grpc-server

![alt text](./doc/img/grpc-web-proxy.png "grpc gateway")

### Protobuf
````
syntax = "proto3";

option go_package = "/helloworld";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a "greeting"
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
````

### Run
```` shell
# 运行proxy和grpc服务
go run .
````


```` shell
# 请求proxy
curl -X POST -H "Content-Type:application/json" -H "Package:/helloworld" \
-H "Service:Greeter" -H "Method:SayHello" \
http://127.0.0.1:8080/ -d '{"name": "Jason"}' 

# proxy返回了grpc response
{"message":"Hello, Jason!"}%   
````
