package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type SelfTest struct {
	UUID       string `json:"id,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	ProblemID  int64  `json:"problem_id,omitempty"`
	Code       string `json:"code,omitempty"`
	Language   string `json:"language,omitempty"`
	Input      string `json:"input,omitempty"`
	IsCompiled bool   `json:"is_compiled,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	TimeCost   int64  `json:"time_cost,omitempty"`
	MemoryCost int64  `json:"memory_cost,omitempty"`
}

type Submission struct {
	UUID        string    `json:"id,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	ProblemID   int64     `json:"problem_id,omitempty"`
	Code        string    `json:"code,omitempty"`
	Language    string    `json:"language,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreateTime  time.Time `json:"create_time"`
	Score       int       `json:"score,omitempty"`
	TimeCost    int64     `json:"time_cost,omitempty"`
	MemoryCost  int64     `json:"memory_cost,omitempty"`
	CaseVersion int8      `json:"case_version,omitempty"`
}

type SubmissionMetadata struct {
	SubmissionID string
	ProblemID    int64
	UserID       int64
	Status       string
	SubmitTime   time.Time
	Score        int32
}

type Case struct {
	Index      int32  `json:"index,omitempty"`
	Score      int32  `json:"score,omitempty"`
	Status     string `json:"status,omitempty"`
	TimeCost   uint64 `json:"time_cost,omitempty"`
	MemoryCost uint64 `json:"memory_cost,omitempty"`
}

type SingleSubmissionDTO struct {
	Submission   *Submission
	CasesResults []*Case `json:"cases_results,omitempty"`
}

type JudgeRepo interface {
	CreateSelfTest(ctx context.Context, selfTest *SelfTest) (string, error)
	CreateSubmission(ctx context.Context, submission *Submission) (string, error)
	GetSubmissions(ctx context.Context, userID int64, problemID int64) ([]*SubmissionMetadata, error)
	GetSingleSubmission(ctx context.Context, submissionID string) (*SingleSubmissionDTO, error)
	GetSelfTest(ctx context.Context, selfTestUUID string) (*SelfTest, error)
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

func (judgeUsecase *JudgeUsecase) GetSingleSubmission(ctx context.Context, submissionID string) (*SingleSubmissionDTO, error) {
	submission, err := judgeUsecase.repo.GetSingleSubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (judgeUsecase JudgeUsecase) GetSelfTest(ctx context.Context, selfTestUUID string) (*SelfTest, error) {
	selfTest, err := judgeUsecase.repo.GetSelfTest(ctx, selfTestUUID)
	if err != nil {
		return nil, err
	}
	return selfTest, nil
}
