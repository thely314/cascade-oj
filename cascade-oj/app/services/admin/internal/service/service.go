package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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
	scheduler           *ContestStatusScheduler
}

func NewAdminService(
	announcementUseCase *biz.AnnouncementUseCase,
	contestUseCase *biz.ContestUsecase,
	logUseCase *biz.LogUseCase,
	problemUseCase *biz.ProblemUsecase,
	submissionUseCase *biz.SubmissionUseCase,
	userUseCase *biz.UserUsecase,
	logger log.Logger,
) *AdminService {
	adminService := &AdminService{
		announcementUseCase: announcementUseCase,
		contestUseCase:      contestUseCase,
		logUseCase:          logUseCase,
		problemUseCase:      problemUseCase,
		submissionUseCase:   submissionUseCase,
		userUseCase:         userUseCase,
	}

	scheduler := NewContestStatusScheduler(contestUseCase, logger, 1*time.Minute)
	scheduler.Start()
	adminService.scheduler = scheduler

	return adminService
}
