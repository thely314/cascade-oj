package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestMetadata struct {
	ID        int64
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

type ContestDetail struct {
	Metadata    ContestMetadata
	Description string
}

type ProblemMetadata struct {
	ID            int64
	Title         string
	ProblemType   string
	TimeLimitMs   int32
	MemoryLimitMb int32
}

type ProblemDetail struct {
	Metadata    ProblemMetadata
	Creator     string
	Description string
}

type ContestRepo interface {
	GetContests(ctx context.Context) ([]*ContestMetadata, error)
	GetSingleContest(ctx context.Context, contestID int64) (*ContestDetail, error)
	JoinContest(ctx context.Context, contestID int64, userID int64) error
	QuitContest(ctx context.Context, contestID int64, userID int64) error
	GetProblems(ctx context.Context, contestID int64) ([]*ProblemMetadata, error)
	GetSingleProblem(ctx context.Context, contestID int64, problemID int64) (*ProblemDetail, error)
}

type ContestUsecase struct {
	repo ContestRepo
	log  *log.Helper
}

func NewContestUsecase(repo ContestRepo, logger log.Logger) *ContestUsecase {
	return &ContestUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (contestUsecase *ContestUsecase) GetContests(ctx context.Context) ([]*ContestMetadata, error) {
	contests, err := contestUsecase.repo.GetContests(ctx)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (contestUsecase *ContestUsecase) GetSingleContest(ctx context.Context, contestID int64) (*ContestDetail, error) {
	contest, err := contestUsecase.repo.GetSingleContest(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (contestUsecase *ContestUsecase) JoinContest(ctx context.Context, contestID int64, userID int64) error {
	err := contestUsecase.repo.JoinContest(ctx, contestID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (contestUsecase *ContestUsecase) QuitContest(ctx context.Context, contestID int64, userID int64) error {
	err := contestUsecase.repo.QuitContest(ctx, contestID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (contestUsecase *ContestUsecase) GetProblems(ctx context.Context, contestID int64) ([]*ProblemMetadata, error) {
	problems, err := contestUsecase.repo.GetProblems(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return problems, nil
}

func (contestUsecase *ContestUsecase) GetSingleProblem(ctx context.Context, contestID int64, problemID int64) (*ProblemDetail, error) {
	problem, err := contestUsecase.repo.GetSingleProblem(ctx, contestID, problemID)
	if err != nil {
		return nil, err
	}
	return problem, nil
}
