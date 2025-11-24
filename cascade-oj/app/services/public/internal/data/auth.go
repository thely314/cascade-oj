package data

import (
	"context"

	"cascade-oj/app/services/public/internal/biz"
	"cascade-oj/ent/user"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
)

type authRepo struct {
	// Data is ent.Client
	data *Data
	log  *log.Helper
}

// find user from database
//
// params:
//   - ctx: context.Context
//   - usernameOrEmail: username or email of the user
//
// returns:
//   - *biz.User: user model
//   - error: nil if success
func (repo *authRepo) FindUserByNameOrEmail(ctx context.Context, usernameOrEmail string) (*biz.User, error) {
	userResult, err := repo.data.db.User.Query().
		Select(user.FieldID, user.FieldUsername, user.FieldEmail, user.FieldPasswordHash, user.FieldRole).
		Where(
			user.Or(
				user.UsernameEQ(usernameOrEmail),
				user.EmailEQ(usernameOrEmail),
			),
		).
		Only(ctx)
	if err != nil {
		repo.log.Errorf("failed to find user by name or email: %v", err)
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

// SignUpAtDatabase creates a new user in the database
//
// params:
//   - ctx: context.Context
//   - username: username of the new user
//   - email: email of the new user
//   - password: password of the new user
//
// returns:
//   - error: nil if success
func (repo *authRepo) SignUpAtDatabase(ctx context.Context, username, email, password string) error {
	bcryptPassword, err := util.GenerateHashPassword(password)
	if err != nil {
		repo.log.Errorf("failed to hash password: %v", err)
		return err
	}
	_, err = repo.data.db.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPasswordHash(bcryptPassword). // In production, hash the password before storing
		SetRole(user.RoleCompetitor).
		Save(ctx)
	if err != nil {
		repo.log.Errorf("failed to create new user: %v", err)
		return err
	}
	return nil
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
