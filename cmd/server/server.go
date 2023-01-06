package main

import (
	"flag"
	"github.com/BingguWang/grpc-gateway-test/cmd/service"
	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
	"google.golang.org/grpc/credentials"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	host = flag.String("host", "localhost", "listening host") // 服务的host
	port = flag.String("port", "50055", "The server port")    // 服务的port
)

/**
普通的grpc协议服务调用
*/
func main() {
	flag.Parsed()
	listen, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		grpclog.Fatalf("监听失败：%v", err)
	}
	log.Printf("server listening at %v", listen.Addr())
	creds, err := credentials.NewServerTLSFromFile(
		"/home/wangbing/grpc-test/key/server.pem",
		"/home/wangbing/grpc-test/key/server.key",
	)
	if err != nil {
		grpclog.Fatalf("Failed to create server TLS credentials %v", err)
	}
	s := grpc.NewServer(
		grpc.Creds(creds),
	)

	pb.RegisterHelloServiceServer(s, &service.HelloServiceImpl{})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
