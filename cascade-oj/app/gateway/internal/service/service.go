package service

import (
	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

type GatewayService struct {
	pb.UnimplementedUserServer
	judgeUsecase   *biz.JudgeUsecase
	contestUsecase *biz.ContestUsecase
	miscUsecase    *biz.MiscUsecase
}

func NewGatewayService(judgeUsecase *biz.JudgeUsecase, contestUsecase *biz.ContestUsecase, miscUsecase *biz.MiscUsecase) *GatewayService {
	return &GatewayService{
		judgeUsecase:   judgeUsecase,
		contestUsecase: contestUsecase,
		miscUsecase:    miscUsecase,
	}
}
