package service

import (
	"context"

	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/broker"
)

func (userService *UserService) ContestCacheConsumer(ctx context.Context, topic string, headers broker.Headers, msg *mq.ContestCacheMsg) error {
	log.Infof("Topic %s, Headers: %+v, Msg: %+v\n", topic, headers, msg)
	err := userService.contestUseCase.ContestCacheDelete(ctx, msg)
	if err != nil {
		log.Errorf("ContestCacheDelete error: %v", err)
	}
	return err
}

func (userService *UserService) ProblemCacheConsumer(ctx context.Context, topic string, headers broker.Headers, msg *mq.ProblemCacheMsg) error {
	log.Infof("Topic %s, Headers: %+v, Msg: %+v\n", topic, headers, msg)
	err := userService.problemUseCase.ProblemCacheDelete(ctx, msg)
	if err != nil {
		log.Errorf("ProblemCacheDelete error: %v", err)
	}
	return err
}
