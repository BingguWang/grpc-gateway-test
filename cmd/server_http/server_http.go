package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	gw "github.com/wbing441282413/grpc-gateway-test/proto/pb"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// grpc服务地址
	addr := "127.0.0.1:50052"

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // 创建一个dialoption数组

	// http转grpc
	err := gw.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, addr, opts)
	if err != nil {
		grpclog.Fatalf("注册处理器错误： ", err)
	}

	http.ListenAndServe(":8080", mux)

}
