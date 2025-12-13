package data

import (
	"context"
	"fmt"

	"cascade-oj/app/services/judge/internal/conf"
	"cascade-oj/app/services/judge/pkg/gojudge"
	"cascade-oj/ent"
	"cascade-oj/ent/problemjudgeconfig"
	filemanage "cascade-oj/pkg/file_manage"
	"cascade-oj/pkg/mq"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	pbGojudge "github.com/criyle/go-judge/pb"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewJudgeRepo)

type Data struct {
	db          *ent.Client
	redis       *redis.Client
	gojudge     *gojudge.GoJudge
	fileManager *filemanage.JudgeConfigManager
	logger      *log.Helper
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	ctx := context.Background()
	logHelper := log.NewHelper(logger)

	// connect to database
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
	if err := client.Schema.Create(ctx); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}

	// check if ProblemJudgeConfig for gojudge engine are created
	isExists, err := client.ProblemJudgeConfig.Query().Where(problemjudgeconfig.JudgeEngineEQ("gojudge")).Exist(ctx)
	if err != nil {
		log.Errorf("failed to query ProblemJudgeConfig: %v", err)
	}
	if !isExists {
		_, err := client.ProblemJudgeConfig.Create().
			SetConfigName("gojudge_default").
			SetDescription("The default config for gojudge").
			SetSubmissionQueueName(mq.GojudgeSubmissionQueueName).
			SetSelfTestQueueName(mq.GojudgeSelfTestQueueName).
			SetJudgeEngine("gojudge").
			Save(ctx)
		if err != nil {
			log.Errorf("failed to create ProblemJudgeConfig: %v", err)
			return nil, nil, err
		}
	}

	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
		DB:   int(c.Redis.Db),
	})
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Errorf("failed to connect to redis: %v", err)
		return nil, nil, err
	}

	// initialize fileManager
	fileManageHelper := filemanage.NewJudgeConfigManager(c.CaseLocation.ProblemCasesLocation)

	// create a gojudge client
	ClientConnection, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(c.JudgeConfig.Endpoint),
		grpc.WithHealthCheck(false),
	)
	if err != nil {
		logHelper.Errorf("failed to create gojudge client: %v", err)
	}
	gojudgeExecClient := pbGojudge.NewExecutorClient(ClientConnection)

	// commands setup
	// extract envs
	envs := make(map[string][]string)
	for k, v := range c.JudgeConfig.Language.Env {
		envs[k] = v.Env
	}
	commands := gojudge.NewCommand(
		c.JudgeConfig.Language.Enable,
		envs,
		c.JudgeConfig.Language.Compile,
		c.JudgeConfig.Language.Run,
		c.JudgeConfig.Language.Source,
		c.JudgeConfig.Language.Target,
		c.JudgeConfig.Language.ExecConfig,
	)

	return &Data{
		db:    client,
		redis: redisClient,
		gojudge: &gojudge.GoJudge{
			Client:   &gojudgeExecClient,
			Commands: &commands,
		},
		fileManager: fileManageHelper,
		logger:      logHelper,
	}, cleanup, nil
}
