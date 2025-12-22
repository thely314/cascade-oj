package biz

import (
	"cascade-oj/ent"
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
	GetProblems(ctx context.Context, contestID int64) ([]*ent.Problem, error)
	GetSingleProblem(ctx context.Context, problemID int64) (*ent.Problem, error)
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
	entProblems, err := problemUsecase.problemRepo.GetProblems(ctx, contestID)
	if err != nil {
		return nil, err
	}
	problems := make([]*Problem, 0, len(entProblems))
	for i := 0; i < len(entProblems); i++ {
		problems = append(problems, &Problem{
			ID:            entProblems[i].ID,
			Title:         entProblems[i].Title,
			TimeLimitMs:   int32(entProblems[i].TimeLimitMs),
			MemoryLimitMb: int32(entProblems[i].MemoryLimitKB),
		})
	}
	return problems, nil
}

func (problemUsecase *ProblemUsecase) GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error) {
	entProblem, err := problemUsecase.problemRepo.GetSingleProblem(ctx, problemID)
	if err != nil {
		return nil, err
	}
	return &DetailedProblem{
		Problem: Problem{
			ID:            entProblem.ID,
			Title:         entProblem.Title,
			TimeLimitMs:   int32(entProblem.TimeLimitMs),
			MemoryLimitMb: int32(entProblem.MemoryLimitKB),
		},
		CreatorUsername: entProblem.Edges.Creator.Username,
		Description:     entProblem.Description,
	}, nil
}
