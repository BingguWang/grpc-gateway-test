package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/wbing441282413/grpc-gateway-test/proto/pb"
)

const addr = "127.0.0.1:50052"

func main() {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatalf("监听失败：%v", err)
	}

	s := grpc.NewServer()

	pb.RegisterHelloServiceServer(s, &HelloServiceImpl{})

	s.Serve(listen)
}

type HelloServiceImpl struct{}

func (m *HelloServiceImpl) Say(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := &pb.HelloResponse{}
	fmt.Println(req.Name)
	if req.Name == "wb" {
		resp.Msg = "调用成功"
	} else {
		resp.Msg = "调用失败"
	}
	return resp, nil
}
