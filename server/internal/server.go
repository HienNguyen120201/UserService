package internal

import (
	"context"
	pb "userservice/server/api"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return nil, nil
}
func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	return nil, nil
}
func (s *server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return nil, nil
}

func RunGRPC(){

}