package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	UserID   int64
	Username string
	Email    string
}

type UserRequestInfo struct {
	Page     int32
	PageSize int32
}

type UserEditInfo struct {
	UserID   int64
	Username string
	Email    string
}

type UserRepo interface {
	GetUsers(ctx context.Context, request UserRequestInfo) ([]*User, error)
	UpdateUserInfo(ctx context.Context, userEditInfo UserEditInfo) (bool, error)
	UpdateUserPassword(ctx context.Context, userID int64, password string) (bool, error)
	DeleteUser(ctx context.Context, userID int64) (bool, error)
}

type UserUsecase struct {
	userRepo UserRepo
	log      *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
		log:      log.NewHelper(logger),
	}
}

func (userUsecase *UserUsecase) GetUsers(ctx context.Context, request UserRequestInfo) ([]*User, error) {
	return userUsecase.userRepo.GetUsers(ctx, request)
}

func (userUsecase *UserUsecase) UpdateUserInfo(ctx context.Context, userEditInfo UserEditInfo) (bool, error) {
	return userUsecase.userRepo.UpdateUserInfo(ctx, userEditInfo)
}

func (userUsecase *UserUsecase) UpdateUserPassword(ctx context.Context, userID int64, password string) (bool, error) {
	return userUsecase.userRepo.UpdateUserPassword(ctx, userID, password)
}

func (userUsecase *UserUsecase) DeleteUser(ctx context.Context, userID int64) (bool, error) {
	return userUsecase.userRepo.DeleteUser(ctx, userID)
}
