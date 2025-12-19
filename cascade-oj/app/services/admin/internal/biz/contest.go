package biz

import (
	"context"
	"time"

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

type ContestCreateInfo struct {
	Title         string
	StartTime     time.Time
	EndTime       time.Time
	ProblemIDList []int64
}

type ContestEditInfo struct {
	ID          int64
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type ContestRank struct {
	UserID   int64
	Username string
	Rank     int32
	Score    int32
}

type ContestRepo interface {
	GetContests(ctx context.Context) ([]*Contest, error)
	GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error)
	PostContest(ctx context.Context, contestCreateInfo ContestCreateInfo) (int64, error)
	PutContest(ctx context.Context, contestEditInfo ContestEditInfo) (bool, error)
	DeleteContest(ctx context.Context, contestID int64) (bool, error)
	GetRanks(ctx context.Context, contestID int64) ([]*ContestRank, error)
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
	return contestUsecase.contestRepo.GetContests(ctx)
}

func (contestUsecase *ContestUsecase) GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error) {
	return contestUsecase.contestRepo.GetSingleContest(ctx, contestID)
}

func (contestUsecase *ContestUsecase) PostContest(ctx context.Context, contestCreateInfo ContestCreateInfo) (int64, error) {
	return contestUsecase.contestRepo.PostContest(ctx, contestCreateInfo)
}
func (contestUsecase *ContestUsecase) PutContest(ctx context.Context, contestEditInfo ContestEditInfo) (bool, error) {
	return contestUsecase.contestRepo.PutContest(ctx, contestEditInfo)
}
func (contestUsecase *ContestUsecase) DeleteContest(ctx context.Context, contestID int64) (bool, error) {
	return contestUsecase.contestRepo.DeleteContest(ctx, contestID)
}
func (contestUsecase *ContestUsecase) GetRanks(ctx context.Context, contestID int64) ([]*ContestRank, error) {
	return contestUsecase.contestRepo.GetRanks(ctx, contestID)
}
