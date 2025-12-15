package server

import (
	"context"

	"cascade-oj/app/services/judge/internal/conf"
	"cascade-oj/app/services/judge/internal/service"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/broker"
	rabbitmqBroker "github.com/tx7do/kratos-transport/broker/rabbitmq"
	"github.com/tx7do/kratos-transport/transport/rabbitmq"
)

func NewMQServer(c *conf.Server, s *service.JudgeService, logger log.Logger) *rabbitmq.Server {
	var opts = []rabbitmq.ServerOption{
		rabbitmq.WithAddress([]string{c.Mq}),
		rabbitmq.WithCodec("json"),
	}

	srv := rabbitmq.NewServer(opts...)

	_ = rabbitmq.RegisterSubscriber(srv, context.Background(), mq.GojudgeSubmissionQueueName, s.SubmissionHandle, broker.WithQueueName(mq.GojudgeSubmissionQueueName), rabbitmqBroker.WithDurableQueue())
	_ = rabbitmq.RegisterSubscriber(srv, context.Background(), mq.GojudgeSelfTestQueueName, s.SelfTestHandle, broker.WithQueueName(mq.GojudgeSelfTestQueueName), rabbitmqBroker.WithDurableQueue())

	return srv
}
