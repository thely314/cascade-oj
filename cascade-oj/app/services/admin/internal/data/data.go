package data

import (
	"context"
	"fmt"

	"cascade-oj/app/services/admin/internal/conf"
	"cascade-oj/ent"
	"cascade-oj/pkg/mq"

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
var ProviderSet = wire.NewSet(NewData, NewAnnouncementRepo, NewContestRepo, NewLogRepo, NewProblemRepo, NewSubmissionRepo, NewUserRepo)

type Data struct {
	db         *ent.Client
	redis      *redis.Client
	mq_channel *mq.MQConnection
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

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
		DB:   int(c.Redis.Db),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Errorf("failed connecting to redis: %v", err)
		return nil, nil, err
	}

	// connect to mq
	mqConn, err := mq.NewMQConnection(c.Mq)
	if err != nil {
		log.Errorf("failed connecting to mq: %v", err)
		return nil, nil, err
	}

	conn, err := mqConn.GetConnection()
	if err != nil {
		log.Errorf("failed opening an init conn")
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("failed opening a channel")
		return nil, nil, err
	}
	defer ch.Close()

	// declare exchanges
	exchangeNames := []string{
		mq.ContestExchangeName,
		mq.ProblemExchangeName,
	}
	for _, exchangeName := range exchangeNames {
		err = mq.NewExchangeDeclare(ch, exchangeName, "topic")
		if err != nil {
			log.Errorf("failed declaring exchange: %v", err)
			return nil, nil, err
		}
	}

	// declare queues
	queueNames := []string{
		mq.ContestCacheQueueName,
		mq.ProblemCacheQueueName,
	}
	for _, queueName := range queueNames {
		_, err = mq.NewQueueDeclare(ch, queueName)
		if err != nil {
			log.Errorf("failed declaring queue: %v", err)
			return nil, nil, err
		}
	}

	// bind queues to exchanges
	// declare bindings as a slice of [queueName, routingKey, exchangeName]
	bindings := [][]string{
		{mq.ContestCacheQueueName, "contest.cache.#", mq.ContestExchangeName},
		{mq.ProblemCacheQueueName, "problem.cache.#", mq.ProblemExchangeName},
	}
	for _, binding := range bindings {
		err = mq.NewBindingDeclare(ch, binding[0], binding[1], binding[2])
		if err != nil {
			log.Errorf("failed declaring binding: %v", err)
			return nil, nil, err
		}
	}

	return &Data{
		db:         client,
		redis:      redisClient,
		mq_channel: mqConn,
	}, cleanup, nil
}

// WithTx 事务封装辅助函数
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) (int64, error)) (int64, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	id, err := fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return 0, rerr
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}
