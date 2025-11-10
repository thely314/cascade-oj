package service

import (
	// TODO api
	"cascade-oj/app/services/public/internal/biz"
)

type AuthService struct {
	// TODO api
	auc *biz.AuthUsecase
}

func NewAuthService(auc *biz.AuthUsecase) *AuthService {
	return &AuthService{
		// TODO api
		auc: auc,
	}
}

// TODO api
// func (as *AuthService) Login((ctx context.Context, req *API_REQ)) (*API_REPLAY, error) {
// jwtToken, err := as.auc.Login(ctx, req.username, req.password)
// if err != nil {
// 	return nil, err
// }
// return &API_REPLAY {
// 	token
// }, nil
// }
