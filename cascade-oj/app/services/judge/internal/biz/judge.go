package biz

import (
	"cascade-oj/pkg/mq"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type JudgeRepo interface {
	JudgeSubmission(ctx context.Context, msg_submission *mq.SubmissionMessage) error
	JudgeSelfTest(ctx context.Context, msg_self_test *mq.SelfTestMessage) error
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

func (judgeUsecase *JudgeUsecase) JudgeSubmission(ctx context.Context, msg *mq.SubmissionMessage) error {
	err := judgeUsecase.repo.JudgeSubmission(ctx, msg)
	if err != nil {
		log.Errorf("JudgeSubmission error: %v", err)
	}
	return err
}

func (judgeUsecase *JudgeUsecase) JudgeSelfTest(ctx context.Context, msg *mq.SelfTestMessage) error {
	err := judgeUsecase.repo.JudgeSelfTest(ctx, msg)
	if err != nil {
		log.Errorf("JudgeSelfTest error: %v", err)
	}
	return err
}
