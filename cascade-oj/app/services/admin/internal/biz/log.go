package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Log struct {
	ID        int64
	Message   string
	TimeStamp time.Time
}

type LogRequestInfo struct {
	Page     int32
	PageSize int32
}

type LogRepo interface {
	GetLogs(ctx context.Context, request LogRequestInfo) ([]*Log, error)
}

type LogUseCase struct {
	repo LogRepo
	log  *log.Helper
}

func NewLogUseCase(repo LogRepo, logger log.Logger) *LogUseCase {
	return &LogUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (ojLogUseCase *LogUseCase) GetLogs(ctx context.Context, request LogRequestInfo) ([]*Log, error) {
	return ojLogUseCase.repo.GetLogs(ctx, request)
}
