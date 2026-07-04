package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"

	"cascade-oj/app/services/judge/internal/conf"
	newlog "cascade-oj/pkg/log"
	"cascade-oj/pkg/metrics"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/tx7do/kratos-transport/transport/rabbitmq"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var metricsHandler http.Handler

// metricsServer is a minimal HTTP server for Prometheus metrics scraping.
var _ transport.Server = (*metricsServer)(nil)

type metricsServer struct{}

func (s *metricsServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler)
	srv := &http.Server{Addr: ":8000", Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *metricsServer) Stop(ctx context.Context) error { return nil }

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "judge"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, ms *rabbitmq.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(ms, &metricsServer{}),
	)
}

func main() {
	flag.Parse()
	logger, err := newlog.NewLogger(
		newlog.WithLevel("INFO"),
		newlog.WithFilename("/data/logs/judge.log"),
		newlog.WithMaxBackups(1),
	)
	if err != nil {
		panic(err)
	}

	log.SetLogger(logger)

	logger = log.With(logger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// Initialize OTel + Prometheus metrics
	var errMetrics error
	metricsHandler, errMetrics = metrics.Init(Name, Version)
	if errMetrics != nil {
		panic(errMetrics)
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		// extra logging for metrics server
		log.Errorf("app.Run failed: %v", err)
		panic(err)
	}
}
