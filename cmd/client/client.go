package main

import (
	"context"
	"flag"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
)

var (
	addr = flag.String("addr", "localhost:50055", "the address to connect to")
)

func main() {
	flag.Parsed()
	creds, err := credentials.NewClientTLSFromFile( // 单向TLS认证
		"/home/wangbing/grpc-test/key/server.pem",
		"x.binggu.example.com",
	)
	clientConn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalf("获取client conn 失败： %v", err)
	}

	client := pb.NewHelloServiceClient(clientConn)
	resp, err := client.Say(context.Background(), &pb.HelloRequest{Name: "wb"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
