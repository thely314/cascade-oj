package data

import (
	"context"
	"fmt"

	"cascade-oj/app/services/user/internal/conf"
	"cascade-oj/ent"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRegisterRepo)

type Data struct {
	db    *ent.Client
	redis *redis.Client
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

	// TODO: connect to redis
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: c.Redis.Addr,
	// 	DB:   int(c.Redis.Db),
	// })
	// if err := redisClient.Ping(context.Background()).Err(); err != nil {
	// 	log.Errorf("failed connecting to redis: %v", err)
	// 	return nil, nil, err
	// }
	// TODO: connect to mq
	// conn, err := amqp.Dial(c.Mq)
	// if err != nil {
	// 	log.Errorf("failed connecting to mq: %v", err)
	// 	return nil, nil, err
	// }
	// ch, err := conn.Channel()
	// if err != nil {
	// 	log.Errorf("failed opening a channel")
	// 	return nil, nil, err
	// }
	// chanMap := make(map[string]*mq.OjChannel)
	// for _, problemType := range client.ProblemType.Query().AllX(context.Background()) {
	// 	sCh, err := mq.NewChannel(ch, problemType.SubmissionChannelName)
	// 	if err != nil {
	// 		log.Errorf("failed creating a submission channel")
	// 		return nil, nil, err
	// 	}
	// 	chanMap[problemType.SubmissionChannelName] = sCh
	// 	stCh, err := mq.NewChannel(ch, problemType.SelfTestChannelName)
	// 	if err != nil {
	// 		log.Errorf("failed creating a self-test channel")
	// 		return nil, nil, err
	// 	}
	// 	chanMap[problemType.SelfTestChannelName] = stCh
	// }
	// return &Data{client, redisClient, chanMap}, cleanup, nil
	// TODO: temporarily return nil redis client
	return &Data{db: client, redis: nil}, cleanup, nil
}
