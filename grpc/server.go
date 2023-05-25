package grpc

import (
	"context"
	"fmt"
	"gateway/codec"
	"gateway/common"
	pb "gateway/grpc/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type GreeterServer interface {
	SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error)
	MustEmbedUnimplementedGreeterServer()
}

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("grpc service SayHello receive name: %s\n", request.Name)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello, %s!", request.Name)}, nil
}

func Start() {
	// 监听127.0.0.1:50051地址
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", common.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc服务端
	encoding.RegisterCodec(codec.JSON{})
	s := grpc.NewServer()

	// 注册Greeter服务
	pb.RegisterGreeterServer(s, &Server{})

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
