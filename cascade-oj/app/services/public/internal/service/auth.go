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

func (as *AuthService) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupReply, error) {
	err := as.auc.Signup(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.SignupReply{
		Code:    200,
		Message: "Signup successful",
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	jwtToken, err := as.auc.Login(ctx, req.UsernameOrEmail, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Code:    200,
		Message: "Login successful",
		Token:   jwtToken,
	}, nil
}

func (as *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	// No server-side logout implementation needed for JWT
	return &pb.LogoutReply{
		Success: true,
	}, nil
}
