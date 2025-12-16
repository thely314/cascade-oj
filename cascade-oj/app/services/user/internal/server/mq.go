package server

import (
	// "context"

	"cascade-oj/app/services/user/internal/conf"
	"cascade-oj/app/services/user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/transport/rabbitmq"
)

// TODO: understanding WIP
const ProblemExchangeName = "update_problem"

func NewMQServer(c *conf.Server, s *service.UserService, logger log.Logger) *rabbitmq.Server {
	var opts = []rabbitmq.ServerOption{
		rabbitmq.WithAddress([]string{c.Mq}),
		rabbitmq.WithCodec("json"),
		rabbitmq.WithExchange(ProblemExchangeName, true),
	}

	srv := rabbitmq.NewServer(opts...)

	// TODO: register mq handlers
	// _ = rabbitmq.RegisterSubscriber(srv, context.Background(), ProblemExchangeName, s.UpdateProblemHandle)

	return srv
}
