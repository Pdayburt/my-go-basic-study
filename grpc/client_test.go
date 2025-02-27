package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"my-go-basic-study/grpc/my-go-basic-study/grpc/userService"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := userService.NewUserServiceClient(conn)
	getById, err := client.GetById(context.Background(), &userService.GetByIdReq{
		Id: 456,
	})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	t.Log(getById.User)
}
