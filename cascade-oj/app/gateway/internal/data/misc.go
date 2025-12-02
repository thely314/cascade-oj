package data

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type miscRepo struct {
	data *Data
	log  *log.Helper
}

func (repo *miscRepo) GetRanks(ctx context.Context, contestID int64) ([]*biz.Rank, error) {
	ranks, err := repo.data.grpcUserClient.GetRanks(ctx, &pb.GetRanksRequest{ContestId: contestID})
	if err != nil {
		return nil, err
	}
	var res []*biz.Rank
	for _, rank := range ranks.Ranks {
		res = append(res, &biz.Rank{
			UserID:   rank.UserId,
			Username: rank.Username,
			UserRank: rank.Rank,
			Score:    rank.Score,
		})
	}
	return res, nil
}

func (repo *miscRepo) GetAnnouncements(ctx context.Context) ([]*biz.Announcement, error) {
	announcements, err := repo.data.grpcUserClient.GetAnnouncements(ctx, &pb.GetAnnouncementsRequest{})
	if err != nil {
		return nil, err
	}
	var res []*biz.Announcement
	for _, announcement := range announcements.Announcements {
		res = append(res, &biz.Announcement{
			ID:            announcement.Id,
			PublisherName: announcement.PublisherName,
			Title:         announcement.Title,
			Content:       announcement.Content,
		})
	}
	return res, nil
}

func (repo *miscRepo) GetUserInfo(ctx context.Context, userID int64) (*biz.UserInfo, error) {
	userInfo, err := repo.data.grpcUserClient.GetUserInfo(ctx, &pb.GetUserInfoRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &biz.UserInfo{
		UserID:   userInfo.UserId,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}, nil
}

func (repo *miscRepo) UpdateUserInfo(ctx context.Context, user *biz.UserInfo) error {
	_, err := repo.data.grpcUserClient.UpdateUserInfo(ctx, &pb.UpdateUserInfoRequest{
		UserId:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewMiscRepo(data *Data, logger log.Logger) biz.MiscRepo {
	return &miscRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
