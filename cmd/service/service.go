package service

import (
	"context"
	"fmt"
	"github.com/BingguWang/grpc-gateway-test/cmd/utils"
	pb "github.com/BingguWang/grpc-gateway-test/proto/mypb"
	"log"
)

type HelloServiceImpl struct{}

func (m *HelloServiceImpl) Say(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("call Service.Say...")
	resp := &pb.HelloResponse{}
	fmt.Println(req.Name)
	if req.Name == "wb" {
		resp.Msg = "it's wb"
	} else {
		resp.Msg = "it's not wb"
	}
	return resp, nil
	// curl -k --cert /home/wangbing/grpc-test/ce-server/server.pem --key /home/wangbing/grpc-test/ce-server/server.key -d '{"name":"wb"}' https://localhost:50055/hello/say
}

func (m *HelloServiceImpl) SayBye(ctx context.Context, req *pb.ByeRequest) (*pb.ByeResponse, error) {
	log.Printf("call Service.SayBye...")
	fmt.Println("req is : ", utils.ToJsonString(req))
	resp := &pb.ByeResponse{Msg: req.Words.String()}
	return resp, nil
	// curl -k --cert /home/wangbing/grpc-test/ce-server/server.pem --key /home/wangbing/grpc-test/ce-server/server.key -d '{"content":"oooop"}'  https://127.0.0.1:50055/hello/bye
}

func (m *HelloServiceImpl) SayGet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("call Service.SayGet...")
	fmt.Println("req is : ", utils.ToJsonString(req))
	resp := &pb.HelloResponse{Msg: req.Name}
	return resp, nil
	// curl -k --cert /home/wangbing/grpc-test/ce-server/server.pem --key /home/wangbing/grpc-test/ce-server/server.key --request GET https://localhost:50055/hello/say/SS
}

func (m *HelloServiceImpl) UpdateHello(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log.Printf("call Service.UpdateHello...")
	fmt.Println("req is : ", req.String())
	resp := &pb.UpdateResponse{Msg: req.Msg}
	return resp, nil
	// curl -k --cert /home/wangbing/grpc-test/ce-server/server.pem --key /home/wangbing/grpc-test/ce-server/server.key --request PATCH -d '{"msg":"fuck"}'  https://127.0.0.1:50055/hello/update
}
