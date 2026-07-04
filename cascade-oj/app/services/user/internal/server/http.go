package server

import (
	"net/http"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/services/user/internal/conf"
	"cascade-oj/app/services/user/internal/service"
	"cascade-oj/pkg/middleware/auth"
	"cascade-oj/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, user *service.UserService, metricsHandler http.Handler, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			metrics.Middleware(),
			// only allow signed-in users to access these apis
			selector.Server(auth.Auth(c.JwtSecret, auth.RoleCompetitor, nil)).
				Path(
					"/api.cascade.user.v1.User/JoinContest",
					"/api.cascade.user.v1.User/QuitContest",
					"/api.cascade.user.v1.User/GetJoinStatus",
					"/api.cascade.user.v1.User/PostSelfTest",
					"/api.cascade.user.v1.User/GetSelfTestResult",
					"/api.cascade.user.v1.User/PostSubmission",
					"/api.cascade.user.v1.User/GetSubmissions",
					"/api.cascade.user.v1.User/GetSingleSubmission",
					"/api.cascade.user.v1.User/GetUserInfo",
					"/api.cascade.user.v1.User/UpdateUserInfo",
				).
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)

	// Register Prometheus /metrics endpoint
	srv.Handle("/metrics", metricsHandler)

	pb.RegisterUserHTTPServer(srv, user)
	return srv
}
