package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Problem struct {
	ID            int64
	Title         string
	TimeLimitMs   int32
	MemoryLimitKB int32
}

type DetailedProblem struct {
	Problem         Problem
	CreatorUsername string
	Description     string
}

type ProblemCreateInfo struct {
	Title           string
	TimeLimitMs     int32
	MemoryLimitKB   int32
	CreatorUsername string
	Description     string
}

type ProblemEditInfo struct {
	ID            int64
	Title         string
	TimeLimitMs   int32
	MemoryLimitKB int32
	Description   string
}

type ProblemRepo interface {
	GetProblems(ctx context.Context, contestID int64) ([]*Problem, error)
	GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error)
	PostProblem(ctx context.Context, problemCreateInfo ProblemCreateInfo) (int64, error)
	PutProblem(ctx context.Context, problemEditInfo ProblemEditInfo) (bool, error)
	DeleteProblem(ctx context.Context, problemID int64) (bool, error)
	PublishProblem(ctx context.Context, problemID int64) (bool, error)
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
	return problemUsecase.problemRepo.GetProblems(ctx, contestID)
}

func (problemUsecase *ProblemUsecase) GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error) {
	return problemUsecase.problemRepo.GetSingleProblem(ctx, problemID)
}

func (problemUsecase *ProblemUsecase) PostProblem(ctx context.Context, problemCreateInfo ProblemCreateInfo) (int64, error) {
	return problemUsecase.problemRepo.PostProblem(ctx, problemCreateInfo)
}
func (problemUsecase *ProblemUsecase) PutProblem(ctx context.Context, problemEditInfo ProblemEditInfo) (bool, error) {
	return problemUsecase.problemRepo.PutProblem(ctx, problemEditInfo)
}
func (problemUsecase *ProblemUsecase) DeleteProblem(ctx context.Context, problemID int64) (bool, error) {
	return problemUsecase.problemRepo.DeleteProblem(ctx, problemID)
}
func (problemUsecase *ProblemUsecase) PublishProblem(ctx context.Context, problemID int64) (bool, error) {
	return problemUsecase.problemRepo.PublishProblem(ctx, problemID)
}
