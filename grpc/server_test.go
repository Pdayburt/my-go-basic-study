package grpc

import (
	"google.golang.org/grpc"
	"my-go-basic-study/grpc/my-go-basic-study/grpc/userService"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	grpcServer := grpc.NewServer()
	userSvc := &Server{}
	userService.RegisterUserServiceServer(grpcServer, userSvc)
	//创建一个监听器 坚挺8090端口
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = grpcServer.Serve(listen)
	if err != nil {
		panic(err)
	}
}
