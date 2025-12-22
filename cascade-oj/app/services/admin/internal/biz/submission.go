package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Submission struct {
	SubmissionUUID string
	ProblemID      int64
	UserID         int64
	Status         string
	SubmitTime     time.Time
	Score          int32
}

type TestCase struct {
	Score      int32
	Status     string
	TimeCost   int32
	MemoryCost int32
}

type DetailedSubmission struct {
	Submission Submission
	Code       string
	Language   string
	TimeCost   int32
	MemoryCost int32
	TestCases  []*TestCase
}

type SubmissionsRequestInfo struct {
	ProblemID int64
	ContestID int64
	UserID    int64
	Page      int32
	PageSize  int32
}

type SubmissionRepo interface {
	GetSubmissions(ctx context.Context, request SubmissionsRequestInfo) ([]*Submission, error)
	GetSingleSubmission(ctx context.Context, submissionUUID string) (*DetailedSubmission, error)
	RejudgeSubmission(ctx context.Context, submissionUUID string) (string, error)
}

type SubmissionUseCase struct {
	submissionRepo SubmissionRepo
	log            *log.Helper
}

func NewSubmissionUseCase(repo SubmissionRepo, logger log.Logger) *SubmissionUseCase {
	return &SubmissionUseCase{
		submissionRepo: repo,
		log:            log.NewHelper(logger),
	}
}

func (submissionUseCase *SubmissionUseCase) GetSubmissions(ctx context.Context, request SubmissionsRequestInfo) ([]*Submission, error) {
	return submissionUseCase.submissionRepo.GetSubmissions(ctx, request)
}

func (submissionUseCase *SubmissionUseCase) GetSingleSubmission(ctx context.Context, submissionUUID string) (*DetailedSubmission, error) {
	return submissionUseCase.submissionRepo.GetSingleSubmission(ctx, submissionUUID)
}

func (submissionUseCase *SubmissionUseCase) RejudgeSubmission(ctx context.Context, submissionUUID string) (string, error) {
	return submissionUseCase.submissionRepo.RejudgeSubmission(ctx, submissionUUID)
}
