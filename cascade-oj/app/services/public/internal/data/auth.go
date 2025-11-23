package data

import (
	"context"

	"cascade-oj/app/services/public/internal/biz"
	"cascade-oj/ent/problemset"
	"cascade-oj/ent/user"

	"github.com/go-kratos/kratos/v2/log"
)

type authRepo struct {
	// Data is ent.Client
	data *Data
	log  *log.Helper
}

// get problem sets from database
func (entry *authRepo) GetProblemSets(ctx context.Context) ([]int64, error) {
	problemSets, err := entry.data.db.ProblemSet.Query().Select(problemset.FieldID).All(ctx)
	if err != nil {
		return nil, err
	}
	var idList []int64
	for _, ps := range problemSets {
		idList = append(idList, ps.ID)
	}
	return idList, nil
}

// find user from database
//
// params:
//   - usernameOrEmail: username or email of the user
//
// returns:
// - *biz.User: user model
func (entry *authRepo) FindUserByNameOrEmail(ctx context.Context, usernameOrEmail string) (*biz.User, error) {
	userResult, err := entry.data.db.User.Query().
		Select(user.FieldID, user.FieldUsername, user.FieldEmail, user.FieldPasswordHash, user.FieldRole).
		Where(
			user.Or(
				user.UsernameEQ(usernameOrEmail),
				user.EmailEQ(usernameOrEmail),
			),
		).
		Only(ctx)
	if err != nil {
		entry.log.Errorf("failed to find user by name or email: %v", err)
		return nil, err
	}
	return &biz.User{
		ID:           userResult.ID,
		Username:     userResult.Username,
		Email:        userResult.Email,
		PasswordHash: userResult.PasswordHash,
		Role:         userResult.Role,
	}, nil
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
