package data

import (
	"context"
	"fmt"

	"cascade-oj/app/services/public/internal/conf"
	"cascade-oj/ent"

	// ent

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
var ProviderSet = wire.NewSet(NewData, NewAuthRepo)

type Data struct {
	db *ent.Client
}

type JWT struct {
	secret     string
	expiration string
	issuer     string
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	driver, err := sql.Open(
		c.Database.Driver,
		c.Database.Source,
	)
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
	if err != nil {
		log.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	return &Data{db: client}, cleanup, nil
}
