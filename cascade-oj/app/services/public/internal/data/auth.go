package data

import (
	"context"

	"cascade-oj/app/services/public/internal/biz"
	// ent

	"github.com/go-kratos/kratos/v2/log"
)

type AuthEntry struct {
	// Data is ent.Client
	data *Data
	log  *log.Helper
}

func (entry *AuthEntry) GetContests(ctx context.Context) ([]int64, error) {
	// TODO
	// get contests from database
	return []int64{}, nil
}

func (entry *AuthEntry) FindUserByNameOrEmail(ctx context.Context, usernameOrEmail string) (*biz.User, error) {
	// TODO
	// find user from database
	return nil, nil
}

func NewAuthEntry(data *Data, logger log.Logger) biz.AuthEntry {
	return &AuthEntry{
		data: data,
		log:  log.NewHelper(logger),
	}
}
