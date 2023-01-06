package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/BingguWang/grpc-gateway-test/cmd/service"
	"github.com/BingguWang/grpc-gateway-test/cmd/utils"
	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	host = flag.String("host", "localhost", "listening host") // 服务的host
	port = flag.String("port", "50055", "The server port")    // 服务的port
)

/**
实现单个服务同时满足grpc调用和http调用
*/

func main() {

	addr := net.JoinHostPort(*host, *port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatalf("TCP Listen err:%v\n", err)
	}

	// grpc tls server
	opts := utils.GetOneSideTlsServerOpts()
	grpcServer := grpc.NewServer(opts...)
	// 注册grpc服务
	pb.RegisterHelloServiceServer(grpcServer, &service.HelloServiceImpl{})

	log.Println("-------开始配置网关，提供http服务")

	// gw server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 因为gateway是要去调grpc server的
	//所以这里gateway相对于grpc server来说是grpc的  客户端
	dopts := utils.GetOneSideTlsClientOpts()
	gwmux := runtime.NewServeMux()
	if err := pb.RegisterHelloServiceHandlerFromEndpoint(ctx, gwmux, addr, dopts); err != nil {
		grpclog.Fatalf("Failed to register gw server: %v\n", err)
	}
	log.Println("-------注册网关成功，提供http服务")

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	srv := &http.Server{
		Addr:      addr,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: getTLSConfig(),
	}
	grpclog.Infof("gRPC and https listen on: %s\n", addr)

	if err := srv.Serve(tls.NewListener(conn, srv.TLSConfig)); err != nil {
		grpclog.Fatal("ListenAndServe: ", err)
	}

	return
}

// 用于判断请求来源于Rpc客户端还是Restful api的请求，根据不同的请求注册不同的ServerHTTP服务
func grpcHandlerFunc(gs *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gs.ServeHTTP(w, r)
		})
	} // TODO 看下http包的使用
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 根据请求头判断是否是grpc调用
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") { // 请求基于HTTP/2且请求的是grpc服务
			gs.ServeHTTP(w, r)
		} else { // 请求的是http服务
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func getTLSConfig() *tls.Config {
	cert, _ := ioutil.ReadFile("/home/wangbing/grpc-test/key/server.pem")
	key, _ := ioutil.ReadFile("/home/wangbing/grpc-test/key/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		grpclog.Fatalf("TLS KeyPair err: %v\n", err)
	}
	demoKeyPair = &pair
	return &tls.Config{
		Certificates: []tls.Certificate{*demoKeyPair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
	}
}

//type HelloServiceImpl struct {
//}
//
//func (m *HelloServiceImpl) Say(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
//	resp := &pb.HelloResponse{}
//	if req.Name == "wb" {
//		resp.Msg = "调用成功"
//	} else {
//		resp.Msg = "调用失败"
//	}
//	return resp, nil
//}
