package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// self test model
type SelfTest struct {
	ID       string
	UserID   int64
	Code     string
	Language string
	Input    string
	Output   string
	Status   string
}

// submission model
type Submission struct {
	ID         string
	UserID     int64
	ProblemID  int64
	Code       string
	Language   string
	Status     string
	CreateTime int64
	Score      int64
	TimeCost   int64
	MemoryCost int64
}

type Case struct {
	Index int32
	State int32
}

type JudgeRepo interface {
	CreateSelfTest(ctx context.Context, selfTest *SelfTest) (string, error)
	CreateSubmission(ctx context.Context, submission *Submission) (string, error)
}

type JudgeUsecase struct {
	repo JudgeRepo
	log  *log.Helper
}

func NewJudgeUsecase(repo JudgeRepo, logger log.Logger) *JudgeUsecase {
	return &JudgeUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (judgeUsecase *JudgeUsecase) CreateSelfTest(ctx context.Context, selfTest *SelfTest) (string, error) {
	// TODO: Implementation of CreateSelfTest
	selfTestID, err := judgeUsecase.repo.CreateSelfTest(ctx, selfTest)
	if err != nil {
		return "", err
	}
	return selfTestID, nil
}

func (judgeUsecase *JudgeUsecase) CreateSubmission(ctx context.Context, submission *Submission) (string, error) {
	// TODO: Implementation of CreateSubmission
	submissionID, err := judgeUsecase.repo.CreateSubmission(ctx, submission)
	if err != nil {
		return "", err
	}
	return submissionID, nil
}
