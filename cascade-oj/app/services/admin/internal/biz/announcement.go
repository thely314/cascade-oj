package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Announcement struct {
	ID            int64
	PublisherName string
	Title         string
	Content       string
}

type AnnouncementCreateInfo struct {
	PublisherID int64
	Title       string
	Content     string
}

type AnnouncementEditInfo struct {
	AnnouncementID int64
	Title          string
	Content        string
}

type AnnouncementRepo interface {
	GetAnnouncements(ctx context.Context) ([]*Announcement, error)
	PostAnnouncement(ctx context.Context, announcementCreateInfo AnnouncementCreateInfo) (int64, error)
	PutAnnouncement(ctx context.Context, announcementEditInfo AnnouncementEditInfo) (bool, error)
	DeleteAnnouncement(ctx context.Context, announcementID int64) (bool, error)
}

type AnnouncementUseCase struct {
	announcementRepo AnnouncementRepo
	log              *log.Helper
}

func NewAnnouncementUseCase(repo AnnouncementRepo, logger log.Logger) *AnnouncementUseCase {
	return &AnnouncementUseCase{
		announcementRepo: repo,
		log:              log.NewHelper(logger),
	}
}

func (announcementUseCase *AnnouncementUseCase) GetAnnouncements(ctx context.Context) ([]*Announcement, error) {
	return announcementUseCase.announcementRepo.GetAnnouncements(ctx)
}

func (announcementUseCase *AnnouncementUseCase) PostAnnouncement(ctx context.Context, announcementCreateInfo AnnouncementCreateInfo) (int64, error) {
	return announcementUseCase.announcementRepo.PostAnnouncement(ctx, announcementCreateInfo)
}

func (announcementUseCase *AnnouncementUseCase) PutAnnouncement(ctx context.Context, announcementEditInfo AnnouncementEditInfo) (bool, error) {
	return announcementUseCase.announcementRepo.PutAnnouncement(ctx, announcementEditInfo)
}

func (announcementUseCase *AnnouncementUseCase) DeleteAnnouncement(ctx context.Context, announcementID int64) (bool, error) {
	return announcementUseCase.announcementRepo.DeleteAnnouncement(ctx, announcementID)
}
