package grpc

import (
	"context"
	"my-go-basic-study/grpc/my-go-basic-study/grpc/userService"
)

type Server struct {
	userService.UnimplementedUserServiceServer // 现在应该可以正确引用了
}

var _ userService.UserServiceServer = (*Server)(nil)

func (s *Server) GetById(ctx context.Context, req *userService.GetByIdReq) (*userService.GetByIdResp, error) {

	return &userService.GetByIdResp{
		User: &userService.User{
			Id:     req.Id,
			Name:   "jack",
			Avatar: "",
		},
	}, nil

}
