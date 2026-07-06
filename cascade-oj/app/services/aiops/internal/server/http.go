package server

import (
	"net/http"

	pb "cascade-oj/api/cascade/aiops/v1"
	"cascade-oj/app/services/aiops/internal/conf"
	"cascade-oj/app/services/aiops/internal/service"
	"cascade-oj/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer creates the AIOps HTTP server.
func NewHTTPServer(
	c *conf.Server,
	svc *service.AIOpsService,
	metricsHandler http.Handler,
	logger log.Logger,
) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			metrics.Middleware(),
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

	pb.RegisterAIOpsHTTPServer(srv, svc)
	return srv
}
