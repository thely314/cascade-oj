package service

import (
	pb "cascade-oj/api/cascade/public/auth/v1"
	"cascade-oj/app/services/public/internal/biz"
	"context"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	auc *biz.AuthUsecase
}

func NewAuthService(auc *biz.AuthUsecase) *AuthService {
	return &AuthService{
		auc: auc,
	}
}

func (as *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	jwtToken, err := as.auc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token: jwtToken,
	}, nil
}
