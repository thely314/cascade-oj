package data

import (
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/announcement"
	"cascade-oj/ent/user"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AnnouncementRepo struct {
	data *Data
	log  *log.Helper
}

func NewAnnouncementRepo(data *Data, logger log.Logger) biz.AnnouncementRepo {
	return &AnnouncementRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (announcementRepo *AnnouncementRepo) GetAnnouncements(ctx context.Context) ([]*biz.Announcement, error) {
	entAnnouncements, err := announcementRepo.data.db.Announcement.
		Query().
		Select(announcement.FieldID, announcement.FieldTitle, announcement.FieldContent).
		WithPublisher(func(uq *ent.UserQuery) {
			uq.Select(
				user.FieldUsername,
			)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}
	announcements := make([]*biz.Announcement, 0, len(entAnnouncements))
	for i := 0; i < len(entAnnouncements); i++ {
		announcements = append(announcements,
			&biz.Announcement{
				ID:            entAnnouncements[i].ID,
				PublisherName: entAnnouncements[i].Edges.Publisher.Username,
				Title:         entAnnouncements[i].Title,
				Content:       entAnnouncements[i].Content,
			},
		)
	}
	return announcements, nil
}
func (announcementRepo *AnnouncementRepo) PostAnnouncement(ctx context.Context, announcementCreateInfo biz.AnnouncementCreateInfo) (int64, error) {
	publisher, err := announcementRepo.data.db.User.
		Query().
		Select(
			user.FieldID,
		).
		Where(
			user.UsernameEQ(announcementCreateInfo.PublisherName),
		).
		Only(ctx)
	if err != nil {
		return -1, err
	}
	newAnnouncementID, err := announcementRepo.data.db.Announcement.Create().SetPublisherID(publisher.ID).SetTitle(announcementCreateInfo.Title).SetContent(announcementCreateInfo.Content).Save(ctx)
	if err != nil {
		return -1, nil
	}
	return newAnnouncementID.ID, err
}
func (announcementRepo *AnnouncementRepo) PutAnnouncement(ctx context.Context, announcementEditInfo biz.AnnouncementEditInfo) (bool, error) {
	err := announcementRepo.data.db.Announcement.
		UpdateOneID(announcementEditInfo.AnnouncementID).
		SetTitle(announcementEditInfo.Title).
		SetContent(announcementEditInfo.Content).
		Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func (announcementRepo *AnnouncementRepo) DeleteAnnouncement(ctx context.Context, announcementID int64) (bool, error) {
	err := announcementRepo.data.db.Announcement.
		DeleteOneID(announcementID).
		Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
