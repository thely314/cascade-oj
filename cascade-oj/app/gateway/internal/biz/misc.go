package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Rank struct {
	UserID   int64
	Username string
	UserRank int32
	Score    int32
}

type Announcement struct {
	ID            int64
	PublisherName string
	Title         string
	Content       string
}

type UserInfo struct {
	UserID   int64
	Username string
	Email    string
}

type MiscRepo interface {
	GetRanks(ctx context.Context, contestID int64) ([]*Rank, error)
	GetAnnouncements(ctx context.Context) ([]*Announcement, error)
	GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, user *UserInfo) error
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

func (miscUsecase *MiscUsecase) GetRanks(ctx context.Context, contestID int64) ([]*Rank, error) {
	ranks, err := miscUsecase.repo.GetRanks(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return ranks, nil
}

func (miscUsecase *MiscUsecase) GetAnnouncements(ctx context.Context) ([]*Announcement, error) {
	announcements, err := miscUsecase.repo.GetAnnouncements(ctx)
	if err != nil {
		return nil, err
	}
	return announcements, nil
}

func (miscUsecase *MiscUsecase) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	userInfo, err := miscUsecase.repo.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (miscUsecase *MiscUsecase) UpdateUserInfo(ctx context.Context, user *UserInfo) error {
	err := miscUsecase.repo.UpdateUserInfo(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
