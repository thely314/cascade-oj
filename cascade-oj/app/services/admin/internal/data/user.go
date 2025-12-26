package data

import (
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent/user"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (userRepo *UserRepo) GetUsers(ctx context.Context, request biz.UserRequestInfo) ([]*biz.User, error) {
	entUsers, err := userRepo.data.db.User.
		Query().
		Select(
			user.FieldID,
			user.FieldUsername,
			user.FieldEmail,
		).
		Offset(int(max(request.Page-1, 0) * request.PageSize)).
		Limit(int(request.PageSize)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*biz.User, 0, len(entUsers))
	for i := 0; i < len(entUsers); i++ {
		users = append(users,
			&biz.User{
				UserID:   entUsers[i].ID,
				Username: entUsers[i].Username,
				Email:    entUsers[i].Email,
			})
	}
	return users, nil
}
func (userRepo *UserRepo) UpdateUserInfo(ctx context.Context, userEditInfo biz.UserEditInfo) (bool, error) {
	err := userRepo.data.db.User.
		UpdateOneID(userEditInfo.UserID).
		SetUsername(userEditInfo.Username).
		SetEmail(userEditInfo.Email).
		Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func (userRepo *UserRepo) DeleteUser(ctx context.Context, userID int64) (bool, error) {
	err := userRepo.data.db.User.DeleteOneID(userID).Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
