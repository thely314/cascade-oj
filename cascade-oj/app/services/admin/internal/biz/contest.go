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
	ProblemIDs  []int64
}

type ContestCreateInfo struct {
	Title         string
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	ProblemIDList []int64
}

type ContestEditInfo struct {
	ID            int64
	Title         string
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	ProblemIDList []int64
}

type ContestRank struct {
	UserID   int64
	Username string
	Rank     int32
	Score    int32
}

type UserInfo struct {
	UserID   int64
	Username string
	Email    string
}

type ContestRepo interface {
	GetContests(ctx context.Context, adminID int64) ([]*Contest, error)
	GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error)
	PostContest(ctx context.Context, contestCreateInfo ContestCreateInfo, adminID int64) (int64, error)
	PutContest(ctx context.Context, contestEditInfo ContestEditInfo) (bool, error)
	DeleteContest(ctx context.Context, contestID int64) (bool, error)
	UpdateContestStatuses(ctx context.Context) error
	GetRanks(ctx context.Context, contestID int64) ([]*ContestRank, error)
	GetContestCompetitors(ctx context.Context, contestID int64) ([]*UserInfo, error)
	PutContestCompetitors(ctx context.Context, contestID int64, userIDs []int64) (bool, error)
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

func (contestUsecase *ContestUsecase) GetContests(ctx context.Context, adminID int64) ([]*Contest, error) {
	return contestUsecase.contestRepo.GetContests(ctx, adminID)
}

func (contestUsecase *ContestUsecase) GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error) {
	return contestUsecase.contestRepo.GetSingleContest(ctx, contestID)
}

func (contestUsecase *ContestUsecase) PostContest(ctx context.Context, contestCreateInfo ContestCreateInfo, adminID int64) (int64, error) {
	return contestUsecase.contestRepo.PostContest(ctx, contestCreateInfo, adminID)
}

func (contestUsecase *ContestUsecase) PutContest(ctx context.Context, contestEditInfo ContestEditInfo) (bool, error) {
	return contestUsecase.contestRepo.PutContest(ctx, contestEditInfo)
}

func (contestUsecase *ContestUsecase) UpdateContestStatuses(ctx context.Context) error {
	return contestUsecase.contestRepo.UpdateContestStatuses(ctx)
}

func (contestUsecase *ContestUsecase) DeleteContest(ctx context.Context, contestID int64) (bool, error) {
	return contestUsecase.contestRepo.DeleteContest(ctx, contestID)
}

func (contestUsecase *ContestUsecase) GetRanks(ctx context.Context, contestID int64) ([]*ContestRank, error) {
	return contestUsecase.contestRepo.GetRanks(ctx, contestID)
}

func (contestUsecase *ContestUsecase) GetContestCompetitors(ctx context.Context, contestID int64) ([]*UserInfo, error) {
	return contestUsecase.contestRepo.GetContestCompetitors(ctx, contestID)
}

func (contestUsecase *ContestUsecase) PutContestCompetitors(ctx context.Context, contestID int64, userIDs []int64) (bool, error) {
	return contestUsecase.contestRepo.PutContestCompetitors(ctx, contestID, userIDs)
}
