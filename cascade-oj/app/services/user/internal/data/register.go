package data

import (
	"context"

	"cascade-oj/app/services/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type registerRepo struct {
	data *Data
	log  *log.Helper
}

func (repo *registerRepo) RegisterGateway(ctx context.Context, register *biz.Register) error {
	// TODO: maybe more implement
	return repo.data.redis.Set(ctx, register.UUID, register.Endpoint, 0).Err()
}

func NewRegisterRepo(data *Data, logger log.Logger) biz.RegisterRepo {
	return &registerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
