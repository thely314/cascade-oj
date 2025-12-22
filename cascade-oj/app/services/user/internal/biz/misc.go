package biz

import (
	"context"

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
	// TODO check user permissions (check invalid userID with JWT)
	return u.repo.GetUserInfoByID(ctx, userID)
}

func (u *MiscUsecase) UpdateUserInfo(ctx context.Context, userID int64, username, email string) error {
	// TODO check user permissions (check invalid userID with JWT)
	return u.repo.UpdateUserInfo(ctx, userID, username, email)
}
