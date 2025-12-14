package service

import (
	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/services/user/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	pb.UnimplementedUserServer
	registerUsecase *biz.RegisterUsecase
	judgeUsecase    *biz.JudgeUsecase
}

func NewUserService(registerUsecase *biz.RegisterUsecase, judgeUsecase *biz.JudgeUsecase) *UserService {
	return &UserService{
		registerUsecase: registerUsecase,
		judgeUsecase:    judgeUsecase,
	}
}
