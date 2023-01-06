package main

import (
	"context"
	"flag"
	"github.com/BingguWang/grpc-gateway-test/cmd/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
	"net"
	"net/http"

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
	// 单向TLS认证，此时反向代理去调gRPC相当于是客户端
	clientOpts := utils.GetOneSideTlsClientOpts()

	mux := runtime.NewServeMux() // 获取一个多路复用器

	// http转grpc, 反向代理去进行带有认证的gRPC调用
	if err := gw.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, addr, clientOpts); err != nil {
		grpclog.Fatalf("注册处理器错误： ", err)
	}
	log.Println("http server listen on 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}

/**
测试http请求
curl -k --cert /home/wangbing/grpc-test/ce-server/server.pem --key /home/wangbing/grpc-test/ce-server/server.key -d '{"name": "john"}' https://localhost:50055
*/
