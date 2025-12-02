package biz

import (
	"context"
	"time"

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
	CreateTime time.Time
	Score      int64
	TimeCost   int64
	MemoryCost int64
}

type SubmissionMetadata struct {
	SubmissionID int64
	ProblemID    int64
	UserID       int64
	Status       string
	SubmitTime   time.Time
	Score        int32
}

type Case struct {
	Index int32
	State int32
}

type JudgeRepo interface {
	CreateSelfTest(ctx context.Context, selfTest *SelfTest) (string, error)
	CreateSubmission(ctx context.Context, submission *Submission) (string, error)
	GetSubmissions(ctx context.Context, userID int64, problemID int64) ([]*SubmissionMetadata, error)
	GetSingleSubmission(ctx context.Context, submissionID string) (*Submission, error)
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
	selfTestID, err := judgeUsecase.repo.CreateSelfTest(ctx, selfTest)
	if err != nil {
		return "", err
	}
	return selfTestID, nil
}

func (judgeUsecase *JudgeUsecase) CreateSubmission(ctx context.Context, submission *Submission) (string, error) {
	submissionID, err := judgeUsecase.repo.CreateSubmission(ctx, submission)
	if err != nil {
		return "", err
	}
	return submissionID, nil
}

func (judgeUsecase *JudgeUsecase) GetSubmissions(ctx context.Context, userID int64, problemID int64) ([]*SubmissionMetadata, error) {
	submissions, err := judgeUsecase.repo.GetSubmissions(ctx, userID, problemID)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (judgeUsecase *JudgeUsecase) GetSingleSubmission(ctx context.Context, submissionID string) (*Submission, error) {
	submission, err := judgeUsecase.repo.GetSingleSubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	return submission, nil
}
