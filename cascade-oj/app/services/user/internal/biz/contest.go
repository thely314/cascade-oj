package biz

import (
	"context"
	"time"

	"cascade-oj/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
)

type Contest struct {
	ID        int64
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

type DetailedContest struct {
	Contest     Contest
	Description string
}

type ContestRepo interface {
	GetContests(ctx context.Context) ([]*Contest, error)
	GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error)
	JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error)
	QuitContest(ctx context.Context, contestID int64, userID int64) (bool, error)
	GetJoinStatus(ctx context.Context, contestID int64, userID int64) (bool, error)
}

type ContestUsecase struct {
	contestRepo ContestRepo
	log         *log.Helper
}

func NewContestUsecase(repo ContestRepo, logger log.Logger) *ContestUsecase {
	return &ContestUsecase{
		contestRepo: repo,
		log:         log.NewHelper(logger),
	}
}

func (contestUsecase *ContestUsecase) GetContests(ctx context.Context) ([]*Contest, error) {
	contests, err := contestUsecase.contestRepo.GetContests(ctx)
	if err != nil {
		return nil, err
	}

	return contests, nil
}

func (contestUsecase *ContestUsecase) GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error) {
	detailedContest, err := contestUsecase.contestRepo.GetSingleContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	return detailedContest, nil
}

func (contestUsecase *ContestUsecase) JoinContest(ctx context.Context, contestID int64) (bool, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	return contestUsecase.contestRepo.JoinContest(ctx, contestID, userID)
}

func (contestUsecase *ContestUsecase) QuitContest(ctx context.Context, contestID int64) (bool, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	return contestUsecase.contestRepo.QuitContest(ctx, contestID, userID)
}

func (contestUsecase *ContestUsecase) GetJoinStatus(ctx context.Context, contestID int64) (bool, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	return contestUsecase.contestRepo.GetJoinStatus(ctx, contestID, userID)
}
