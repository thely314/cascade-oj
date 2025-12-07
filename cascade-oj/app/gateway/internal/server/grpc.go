package server

import (
	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/conf"
	"cascade-oj/app/gateway/internal/service"
	"cascade-oj/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(c *conf.Server, gateway *service.GatewayService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			auth.Auth(c.JwtSecret, auth.RoleCompetitor, nil),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	server := grpc.NewServer(opts...)
	pb.RegisterUserServer(server, gateway)
	return server
}
