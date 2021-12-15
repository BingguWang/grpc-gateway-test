package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/wbing441282413/grpc-gateway-test/proto/pb"
)

const addr = "127.0.0.1:50052"

func main() {

	endpoint := "127.0.0.1:50052"
	conn, err := net.Listen("tcp", endpoint)
	if err != nil {
		grpclog.Fatalf("TCP Listen err:%v\n", err)
	}

	// grpc tls server
	creds, err := credentials.NewServerTLSFromFile("../../keys/server.pem", "../../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to create server TLS credentials %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServiceServer(grpcServer, &HelloServiceImpl{})

	// gw server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 因为gateway是要去调grpc server的，所以这里gateway相对于grpc server来说是grpc的客户端
	dcreds, err := credentials.NewClientTLSFromFile("../../keys/server.pem", "localhost")
	if err != nil {
		grpclog.Fatalf("Failed to create client TLS credentials %v", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	if err = pb.RegisterHelloServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil {
		grpclog.Fatalf("Failed to register gw server: %v\n", err)
	}

	// http服务
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	srv := &http.Server{
		Addr:      endpoint,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: getTLSConfig(),
	}

	grpclog.Infof("gRPC and https listen on: %s\n", endpoint)

	if err = srv.Serve(tls.NewListener(conn, srv.TLSConfig)); err != nil {
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
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") { // 请求基于HTTP/2且请求的是grpc服务
			gs.ServeHTTP(w, r)
		} else { // 请求的是http服务
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func getTLSConfig() *tls.Config {
	cert, _ := ioutil.ReadFile("../../keys/server.pem")
	key, _ := ioutil.ReadFile("../../keys/server.key")
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

type HelloServiceImpl struct {
}

func (m *HelloServiceImpl) Say(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := &pb.HelloResponse{}
	if req.Name == "wb" {
		resp.Msg = "调用成功"
	} else {
		resp.Msg = "调用失败"
	}
	return resp, nil
}
