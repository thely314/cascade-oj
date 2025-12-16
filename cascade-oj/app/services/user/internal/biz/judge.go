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
	CompileMsg string `json:"compile_msg,omitempty"`
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
	CompileMsg  string    `json:"compile_msg,omitempty"`
}

type Case struct {
	Index      int32  `json:"index,omitempty"`
	Score      int32  `json:"score,omitempty"`
	Status     int32  `json:"status,omitempty"`
	TimeCost   uint64 `json:"time_cost,omitempty"`
	MemoryCost uint64 `json:"memory_cost,omitempty"`
}

type JudgeRepo interface {
	CreateSelfTest(ctx context.Context, selfTest *SelfTest) (string, error)
	CreateSubmission(ctx context.Context, submission *Submission) (string, error)
	GetSelfTest(ctx context.Context, submissionID string) (*SelfTest, error)
	GetSubmissions(ctx context.Context, contestID int64, problemID int64) ([]*Submission, error)
	GetSingleSubmission(ctx context.Context, submissionID string) (*Submission, error)
	GetCases(ctx context.Context, submissionID string) ([]*Case, error)
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
	return judgeUsecase.repo.CreateSelfTest(ctx, selfTest)
}

func (judgeUsecase *JudgeUsecase) CreateSubmission(ctx context.Context, submission *Submission) (string, error) {
	return judgeUsecase.repo.CreateSubmission(ctx, submission)
}

func (judgeUsecase *JudgeUsecase) GetSelfTest(ctx context.Context, submissionID string) (*SelfTest, error) {
	return judgeUsecase.repo.GetSelfTest(ctx, submissionID)
}

func (judgeUsecase *JudgeUsecase) GetSubmissions(ctx context.Context, contestID int64, problemID int64) ([]*Submission, error) {
	return judgeUsecase.repo.GetSubmissions(ctx, contestID, problemID)
}

func (judgeUsecase *JudgeUsecase) GetSingleSubmission(ctx context.Context, submissionID string) (*Submission, error) {
	return judgeUsecase.repo.GetSingleSubmission(ctx, submissionID)
}

func (judgeUsecase *JudgeUsecase) GetCases(ctx context.Context, submissionID string) ([]*Case, error) {
	return judgeUsecase.repo.GetCases(ctx, submissionID)
}
