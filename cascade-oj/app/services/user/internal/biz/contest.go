package biz

import (
	"cascade-oj/ent"
	"context"
	"errors"
	"fmt"
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
	GetContests(ctx context.Context) ([]*ent.ProblemSet, error)
	GetSingleContest(ctx context.Context, contestID int64) (*ent.ProblemSet, error)
	JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error)
	QuitContest(ctx context.Context, contestID int64, userID int64) (bool, error)
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
	entContests, err := contestUsecase.contestRepo.GetContests(ctx)
	if err != nil {
		return nil, err
	}

	contests := make([]*Contest, 0, len(entContests))
	for i := 0; i < len(entContests); i++ {
		contests = append(contests, &Contest{
			ID:        entContests[i].ID,
			Title:     entContests[i].Name,
			StartTime: entContests[i].StartTime,
			EndTime:   entContests[i].EndTime,
			Status:    string(entContests[i].Status),
		})
	}

	return contests, nil
}

func (contestUsecase *ContestUsecase) GetSingleContest(ctx context.Context, contestID int64) (*DetailedContest, error) {
	entContest, err := contestUsecase.contestRepo.GetSingleContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	return &DetailedContest{
		Contest: Contest{
			ID:        entContest.ID,
			Title:     entContest.Name,
			StartTime: entContest.StartTime,
			EndTime:   entContest.EndTime,
			Status:    string(entContest.Status),
		},
		Description: entContest.Description,
	}, nil
}

func (contestUsecase *ContestUsecase) JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	if userID != ctx.Value("userInfo").(*auth.Claims).UserID {
		return false, errors.New(fmt.Sprintf("user %d trying to join contest as user %d", ctx.Value("userInfo").(*auth.Claims).UserID, userID))
	}
	return contestUsecase.contestRepo.JoinContest(ctx, contestID, userID)
}

func (contestUsecase *ContestUsecase) QuitContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	if userID != ctx.Value("userInfo").(*auth.Claims).UserID {
		return false, errors.New(fmt.Sprintf("user %d trying to quit contest as user %d", ctx.Value("userInfo").(*auth.Claims).UserID, userID))
	}
	return contestUsecase.contestRepo.QuitContest(ctx, contestID, userID)
}
