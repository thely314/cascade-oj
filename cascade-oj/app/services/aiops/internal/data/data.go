package data

import (
	"context"
	"fmt"

	"cascade-oj/app/services/aiops/internal/conf"
	"cascade-oj/ent"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAlertEventRepo, NewAlertReportRepo)

// Data holds the database client.
type Data struct {
	db *ent.Client
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	driver, err := sql.Open(
		c.Database.Driver,
		c.Database.Source,
	)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed opening connection to mysql: %w", err)
	}
	sqlDriver := dialect.DebugWithContext(driver, func(ctx context.Context, i ...interface{}) {
		log.WithContext(ctx, logger)
		tracer := otel.Tracer("ent.")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
	client := ent.NewClient(ent.Driver(sqlDriver))
	return &Data{db: client}, cleanup, nil
}

func (d *Data) DB() *ent.Client {
	return d.db
}
