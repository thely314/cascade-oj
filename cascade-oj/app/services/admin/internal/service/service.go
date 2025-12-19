package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminService)

type AdminService struct {
	pb.UnimplementedAdminServer
	announcementUseCase *biz.AnnouncementUseCase
	contestUseCase      *biz.ContestUsecase
	logUseCase          *biz.LogUseCase
	problemUseCase      *biz.ProblemUsecase
	submissionUseCase   *biz.SubmissionUseCase
	userUseCase         *biz.UserUsecase
}

func NewAdminService(
	announcementUseCase *biz.AnnouncementUseCase,
	contestUseCase *biz.ContestUsecase,
	logUseCase *biz.LogUseCase,
	problemUseCase *biz.ProblemUsecase,
	submissionUseCase *biz.SubmissionUseCase,
	userUseCase *biz.UserUsecase,
) *AdminService {
	return &AdminService{
		announcementUseCase: announcementUseCase,
		contestUseCase:      contestUseCase,
		logUseCase:          logUseCase,
		problemUseCase:      problemUseCase,
		submissionUseCase:   submissionUseCase,
		userUseCase:         userUseCase,
	}
}
