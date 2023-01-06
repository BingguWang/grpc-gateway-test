package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BingguWang/grpc-gateway-test/cmd/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
)

var (
	addr = flag.String("addr", "localhost:50055", "the address to connect to")
)

func main() {
	flag.Parsed()
	// 单向TLS认证
	opts := utils.GetOneSideTlsClientOpts()
	clientConn, err := grpc.Dial(*addr, opts...)
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
