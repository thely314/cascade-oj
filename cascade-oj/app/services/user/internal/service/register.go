package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/services/user/internal/biz"

	"github.com/google/uuid"
)

func (userService *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	ret := uuid.New().String()
	err := userService.registerUsecase.RegisterGateway(ctx, &biz.Register{
		UUID:     ret,
		Endpoint: req.Endpoint,
		Secret:   req.ApiKey,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{
		Token: ret,
	}, nil
}
