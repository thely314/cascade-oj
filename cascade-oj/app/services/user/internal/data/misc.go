package data

import (
	"context"

	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/problemset_user"
	"cascade-oj/ent/user"

	"github.com/go-kratos/kratos/v2/log"
)

type miscRepo struct {
	data *Data
	log  *log.Helper
}

func NewMiscRepo(data *Data, logger log.Logger) biz.MiscRepo {
	return &miscRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *miscRepo) GetRanks(ctx context.Context, contestID int64) ([]*biz.Rank, error) {
	ranks, err := repo.data.db.ProblemSet_User.Query().
		Where(problemset_user.ProblemSetIDEQ(contestID)).
		Order(ent.Desc("score")).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var res []*biz.Rank
	for i, rank := range ranks {
		res = append(res, &biz.Rank{
			UserID:   rank.UserID,
			Username: rank.Edges.User.Username,
			Score:    int64(rank.TotalScore),
			Rank:     int32(i + 1),
		})
	}
	return res, nil
}

func (repo *miscRepo) GetAnnouncements(ctx context.Context) ([]*biz.Announcement, error) {
	announcements, err := repo.data.db.Announcement.Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	var res []*biz.Announcement
	for _, announcement := range announcements {
		res = append(res, &biz.Announcement{
			ID:             announcement.ID,
			Publisher_name: announcement.Edges.Publisher.Username,
			Title:          announcement.Title,
			Content:        announcement.Content,
		})
	}
	return res, nil
}

func (repo *miscRepo) GetUserInfoByID(ctx context.Context, userID int64) (*biz.UserInfo, error) {
	user, err := repo.data.db.User.Query().
		Where(user.IDEQ(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.UserInfo{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (repo *miscRepo) UpdateUserInfo(ctx context.Context, userID int64, username, email string) error {
	_, err := repo.data.db.User.Update().
		Where(user.IDEQ(userID)).
		SetUsername(username).
		SetEmail(email).
		Save(ctx)
	return err
}
