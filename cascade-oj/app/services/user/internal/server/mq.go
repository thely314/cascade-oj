package server

import (
	"context"

	"cascade-oj/app/services/user/internal/conf"
	"cascade-oj/app/services/user/internal/service"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/broker"
	rabbitmqBroker "github.com/tx7do/kratos-transport/broker/rabbitmq"
	"github.com/tx7do/kratos-transport/transport/rabbitmq"
)

func NewMQServer(c *conf.Server, s *service.UserService, logger log.Logger) *rabbitmq.Server {
	var opts = []rabbitmq.ServerOption{
		rabbitmq.WithAddress([]string{c.Mq}),
		rabbitmq.WithCodec("json"),
	}

	srv := rabbitmq.NewServer(opts...)

	// TODO kratos-transport 没有支持 exchange 与 queue 的绑定
	// 当前 server 实现为通过默认 direct exchange 传递消息，支持定义 handler 处理
	// 如果使用非默认 exchange，需要手动接收 consume Channel 并处理
	// 暂时先使用默认 exchange，快速开发
	_ = rabbitmq.RegisterSubscriber(srv,
		context.Background(),
		mq.ContestCacheQueueName,
		s.ContestCacheConsumer,
		broker.WithQueueName(mq.ContestCacheQueueName),
		rabbitmqBroker.WithDurableQueue(),
	)
	_ = rabbitmq.RegisterSubscriber(srv,
		context.Background(),
		mq.ProblemCacheQueueName,
		s.ProblemCacheConsumer,
		broker.WithQueueName(mq.ProblemCacheQueueName),
		rabbitmqBroker.WithDurableQueue(),
	)

	return srv
}
