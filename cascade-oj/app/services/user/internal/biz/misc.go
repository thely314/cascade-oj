package biz

import (
	"context"
	"errors"
	"fmt"

	"cascade-oj/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
)

type Rank struct {
	UserID   int64
	Username string
	Score    int64
	Rank     int32
}

type Announcement struct {
	ID             int64
	Publisher_name string
	Title          string
	Content        string
}

type UserInfo struct {
	UserID   int64
	Username string
	Email    string
}

type MiscRepo interface {
	GetRanks(ctx context.Context, contestID int64) ([]*Rank, error)
	GetAnnouncements(ctx context.Context) ([]*Announcement, error)
	GetUserInfoByID(ctx context.Context, userID int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, userID int64, username, email string) error
}

type MiscUsecase struct {
	repo MiscRepo
	log  *log.Helper
}

func NewMiscUsecase(repo MiscRepo, logger log.Logger) *MiscUsecase {
	return &MiscUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (u *MiscUsecase) GetRanks(ctx context.Context, contestID int64) ([]*Rank, error) {
	return u.repo.GetRanks(ctx, contestID)
}

func (u *MiscUsecase) GetAnnouncements(ctx context.Context) ([]*Announcement, error) {
	return u.repo.GetAnnouncements(ctx)
}

func (u *MiscUsecase) GetUserInfoByID(ctx context.Context, userID int64) (*UserInfo, error) {
	if userID != ctx.Value("userInfo").(*auth.Claims).UserID {
		return nil, errors.New(fmt.Sprintf("user %d trying to get user info as user %d", ctx.Value("userInfo").(*auth.Claims).UserID, userID))
	}
	return u.repo.GetUserInfoByID(ctx, userID)
}

func (u *MiscUsecase) UpdateUserInfo(ctx context.Context, userID int64, username, email string) error {
	if userID != ctx.Value("userInfo").(*auth.Claims).UserID {
		return errors.New(fmt.Sprintf("user %d trying to update user info as user %d", ctx.Value("userInfo").(*auth.Claims).UserID, userID))
	}
	return u.repo.UpdateUserInfo(ctx, userID, username, email)
}
