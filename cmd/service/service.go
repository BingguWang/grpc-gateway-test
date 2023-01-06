package service

import (
	"context"
	"fmt"
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
}
