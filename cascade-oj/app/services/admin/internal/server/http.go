package server

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/conf"
	"cascade-oj/app/services/admin/internal/service"
	"cascade-oj/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, admin *service.AdminService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			auth.Auth(c.JwtSecret, auth.RoleAdmin, nil),
		),
		// http.Filter(handlers.CORS(
		// 	handlers.AllowedOrigins([]string{"*"}),
		// 	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		// 	handlers.AllowedHeaders([]string{"Content-Type", "token", "Token"}),
		// )),
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
	srv := http.NewServer(opts...)
	pb.RegisterAdminHTTPServer(srv, admin)
	return srv
}
