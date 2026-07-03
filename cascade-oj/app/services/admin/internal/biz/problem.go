package biz

import (
	"context"
	"errors"
	"io"

	"github.com/go-kratos/kratos/v2/log"
)

type ProblemStatus int32

const (
	ProblemStatusUnavailable ProblemStatus = 0
	ProblemStatusAvailable   ProblemStatus = 1
	ProblemStatusUsing       ProblemStatus = 2
	ProblemStatusDeleted     ProblemStatus = 3
)

type ProblemTemplate struct {
	Name    string
	Content string
}

type Problem struct {
	ID            int64
	Title         string
	TimeLimitMs   int32
	MemoryLimitKB int32
	Status        ProblemStatus
	Templates     []*ProblemTemplate
	Description   string
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
	Templates       []*ProblemTemplate
}

type ProblemEditInfo struct {
	ID            int64
	Title         string
	TimeLimitMs   int32
	MemoryLimitKB int32
	Description   string
	Templates     []*ProblemTemplate
}

type ProblemRepo interface {
	GetProblems(ctx context.Context, contestID int64) ([]*Problem, error)
	GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error)
	PostProblem(ctx context.Context, creatorID int64, info ProblemCreateInfo) (int64, error)
	PutProblem(ctx context.Context, info ProblemEditInfo) (bool, error)
	DeleteProblem(ctx context.Context, problemID int64) (bool, error)
	PublishProblem(ctx context.Context, problemID int64) (bool, error)
	DisableProblem(ctx context.Context, problemID int64) (bool, error)
	SaveTestCases(ctx context.Context, problemID int64, file io.Reader) error
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

func (uc *ProblemUsecase) HandleTestCasesUpload(ctx context.Context, problemID int64, file io.Reader, filename string) error {
	return uc.problemRepo.SaveTestCases(ctx, problemID, file)
}

func (problemUsecase *ProblemUsecase) GetProblems(ctx context.Context, contestID int64) ([]*Problem, error) {
	return problemUsecase.problemRepo.GetProblems(ctx, contestID)
}

func (problemUsecase *ProblemUsecase) GetSingleProblem(ctx context.Context, problemID int64) (*DetailedProblem, error) {
	return problemUsecase.problemRepo.GetSingleProblem(ctx, problemID)
}

func (problemUsecase *ProblemUsecase) PostProblem(ctx context.Context, creatorID int64, problemCreateInfo ProblemCreateInfo) (int64, error) {
	return problemUsecase.problemRepo.PostProblem(ctx, creatorID, problemCreateInfo)
}

func (problemUsecase *ProblemUsecase) PutProblem(ctx context.Context, problemEditInfo ProblemEditInfo) (bool, error) {
	detailedProblem, err := problemUsecase.problemRepo.GetSingleProblem(ctx, problemEditInfo.ID)
	if err != nil {
		return false, err
	}

	if detailedProblem.Problem.Status == ProblemStatusUsing || detailedProblem.Problem.Status == ProblemStatusDeleted {
		return false, errors.New("business rule error: problems in use or deleted cannot be edited")
	}

	return problemUsecase.problemRepo.PutProblem(ctx, problemEditInfo)
}

func (problemUsecase *ProblemUsecase) DeleteProblem(ctx context.Context, problemID int64) (bool, error) {
	return problemUsecase.problemRepo.DeleteProblem(ctx, problemID)
}

func (problemUsecase *ProblemUsecase) PublishProblem(ctx context.Context, problemID int64) (bool, error) {
	return problemUsecase.problemRepo.PublishProblem(ctx, problemID)
}

func (problemUsecase *ProblemUsecase) DisableProblem(ctx context.Context, problemID int64) (bool, error) {
	return problemUsecase.problemRepo.DisableProblem(ctx, problemID)
}
