package server

import (
	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/conf"
	"cascade-oj/app/gateway/internal/service"
	"cascade-oj/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *conf.Server, gateway *service.GatewayService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			auth.Auth(c.JwtSecret, auth.RoleCompetitor, nil),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	server := http.NewServer(opts...)
	pb.RegisterUserHTTPServer(server, gateway)
	return server
}
