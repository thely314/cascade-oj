// Package metrics provides shared OpenTelemetry + Prometheus metrics setup
// for all Cascade OJ microservices.
package metrics

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	kmiddleware "github.com/go-kratos/kratos/v2/middleware"
	kmetrics "github.com/go-kratos/kratos/v2/middleware/metrics"
)

var (
	// Requests is the default requests counter.
	Requests metric.Int64Counter
	// Seconds is the default seconds histogram.
	Seconds metric.Float64Histogram
)

// Init sets up the OTel meter provider with a Prometheus exporter.
// It returns an http.Handler for the /metrics endpoint.
func Init(serviceName, serviceVersion string) (http.Handler, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)

	meter := otel.Meter(serviceName,
		metric.WithInstrumentationVersion(serviceVersion),
	)

	Requests, err = kmetrics.DefaultRequestsCounter(meter, kmetrics.DefaultServerRequestsCounterName)
	if err != nil {
		return nil, err
	}

	Seconds, err = kmetrics.DefaultSecondsHistogram(meter, kmetrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}

	// promhttp.Handler() serves metrics from the default registry,
	// which the OTel Prometheus exporter writes to.
	return promhttp.Handler(), nil
}

// Middleware returns a pre-configured metrics middleware for HTTP/gRPC servers.
func Middleware() kmiddleware.Middleware {
	return kmetrics.Server(
		kmetrics.WithSeconds(Seconds),
		kmetrics.WithRequests(Requests),
	)
}
