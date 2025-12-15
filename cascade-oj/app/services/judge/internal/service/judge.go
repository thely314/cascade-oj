package service

import (
	"context"

	"cascade-oj/app/services/judge/internal/biz"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/broker"
)

type JudgeService struct {
	judgeUsecase *biz.JudgeUsecase
}

func NewJudgeService(judgeUsecase *biz.JudgeUsecase) *JudgeService {
	return &JudgeService{
		judgeUsecase: judgeUsecase,
	}
}

func (judgeService *JudgeService) SubmissionHandle(ctx context.Context, topic string, headers broker.Headers, msg *mq.SubmissionMessage) error {
	log.Infof("Topic %s, Headers: %+v, Msg: %+v\n", topic, headers, msg)
	err := judgeService.judgeUsecase.JudgeSubmission(ctx, msg)
	if err != nil {
		log.Errorf("JudgeSubmission error: %v", err)
	}
	return err
}

func (judgeService *JudgeService) SelfTestHandle(ctx context.Context, topic string, headers broker.Headers, msg *mq.SelfTestMessage) error {
	log.Infof("Topic %s, Headers: %+v, Msg: %+v\n", topic, headers, msg)
	err := judgeService.judgeUsecase.JudgeSelfTest(ctx, msg)
	if err != nil {
		log.Errorf("JudgeSelfTest error: %v", err)
	}
	return err
}
