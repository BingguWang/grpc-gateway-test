package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	gw "github.com/BingguWang/grpc-gateway-test/proto/mypb"
)

var (
	host = flag.String("host", "localhost", "listening host") // 服务的host
	port = flag.String("port", "50055", "The server port")    // 服务的port
)

/**
通过gateway生成的反向代理将http请求转为内部的grpc协议调用，最终将服务调用结果转为满足http协议的响应并返回
*/
func main() {
	flag.Parsed()
	// grpc服务地址
	addr := net.JoinHostPort(*host, *port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	creds, err := credentials.NewClientTLSFromFile( // 单向TLS认证
		"/home/wangbing/grpc-test/key/server.pem",
		"x.binggu.example.com",
	)
	if err != nil {
		grpclog.Fatalf("获取client conn 失败： %v", err)
	}
	mux := runtime.NewServeMux() // 获取一个多路复用器
	opts := []grpc.DialOption{
		//grpc.WithTransportCredentials(insecure.NewCredentials()), // 无需认证
		grpc.WithTransportCredentials(creds),
	}

	// http转grpc, 反向代理去进行带有认证的gRPC调用
	if err := gw.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		grpclog.Fatalf("注册处理器错误： ", err)
	}
	log.Println("http server listen on 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
