package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/wbing441282413/grpc-gateway-test/proto/pb"
)

const addr = "127.0.0.1:50052"

func main() {

	creds, err := credentials.NewClientTLSFromFile("../../keys/server.pem", "localhost")
	clientConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalf("获取client conn 失败： %v", err)
	}

	client := pb.NewHelloServiceClient(clientConn)
	resp, err := client.Say(context.Background(), &pb.HelloRequest{Name: "wb"})
	fmt.Println(resp)
}
