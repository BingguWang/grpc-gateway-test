package main

import (
	"flag"
	"github.com/BingguWang/grpc-gateway-test/cmd/service"
	"github.com/BingguWang/grpc-gateway-test/cmd/utils"
	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
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
	// 单向TLS认证
	opts := utils.GetOneSideTlsServerOpts()
	s := grpc.NewServer(opts...)

	pb.RegisterHelloServiceServer(s, &service.HelloServiceImpl{})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
