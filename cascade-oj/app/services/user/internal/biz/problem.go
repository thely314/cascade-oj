package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Problem struct {
	ID            int64
	Title         string
	TimeLimitMs   int32
	MemoryLimitMb int32
}

type DetailedProblem struct {
	Problem         Problem
	CreatorUsername string
	Description     string
}

type ProblemRepo interface {
	GetProblems(ctx context.Context, contestID int64) ([]*Problem, error)
	GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error)
}

type ProblemUsecase struct {
	problemRepo ProblemRepo
	log         *log.Helper
}

func NewProblemUsecase(repo ProblemRepo, logger log.Logger) *ProblemUsecase {
	return &ProblemUsecase{
		problemRepo: repo,
		log:         log.NewHelper(logger),
	}
}

func (problemUsecase *ProblemUsecase) GetProblems(ctx context.Context, contestID int64) ([]*Problem, error) {
	problems, err := problemUsecase.problemRepo.GetProblems(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return problems, nil
}

func (problemUsecase *ProblemUsecase) GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error) {
	detailedProblem, err := problemUsecase.problemRepo.GetSingleProblem(ctx, problemID)
	if err != nil {
		return nil, err
	}
	return detailedProblem, nil
}
