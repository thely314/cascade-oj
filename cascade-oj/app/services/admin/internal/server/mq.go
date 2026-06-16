package server

import (
	// "context"

	"cascade-oj/app/services/admin/internal/conf"
	"cascade-oj/app/services/admin/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/transport/rabbitmq"
)

func NewMQServer(c *conf.Server, s *service.AdminService, logger log.Logger) *rabbitmq.Server {
	var opts = []rabbitmq.ServerOption{
		rabbitmq.WithAddress([]string{c.Mq}),
		rabbitmq.WithCodec("json"),
	}

	srv := rabbitmq.NewServer(opts...)

	return srv
}
