package biz

import (
	"context"
	"errors"
	"strconv"
	"time"

	"cascade-oj/app/services/public/internal/conf"
	"cascade-oj/ent/user"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

// user model
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Role         user.Role // "competitor", "creator", "admin"
}

type MyCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthRepo is a entry of searching user.
type AuthRepo interface {
	FindUserByNameOrEmail(context.Context, string) (*User, error)
	SignUpAtDatabase(context.Context, string, string, string) error
}

type JwtConfig struct {
	secret     string
	expiration time.Duration
	issuer     string
}

// AuthUsecase is a Auth usecase.
type AuthUsecase struct {
	authRepo AuthRepo
	log      *log.Helper
	jwtConf  JwtConfig
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger, jwtConfig *conf.Jwt) *AuthUsecase {
	var expiryTime int32 = 24 // default 24 hours
	if jwtConfig.Expiration != 0 {
		expiryTime = jwtConfig.Expiration
	}
	var issuer string = "cascade-oj"
	if jwtConfig.Issuer != "" {
		issuer = jwtConfig.Issuer
	}
	return &AuthUsecase{
		authRepo: repo,
		log:      log.NewHelper(logger),
		jwtConf: JwtConfig{
			secret:     jwtConfig.Secret,
			expiration: time.Duration(expiryTime) * time.Hour,
			issuer:     issuer,
		},
	}
}

func (authUsecase *AuthUsecase) Signup(ctx context.Context, username string, email string, password string) error {
	userWithSameName, err := authUsecase.authRepo.FindUserByNameOrEmail(ctx, username)
	if err == nil && userWithSameName != nil {
		return errors.New("username already exists")
	}
	userWithSameEmail, err := authUsecase.authRepo.FindUserByNameOrEmail(ctx, email)
	if err == nil && userWithSameEmail != nil {
		return errors.New("email already registered")
	}
	err = authUsecase.authRepo.SignUpAtDatabase(ctx, username, email, password)
	if err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func (authUsecase *AuthUsecase) Login(ctx context.Context, usernameOrEmail string, password string) (string, error) {
	user, err := authUsecase.authRepo.FindUserByNameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return "", errors.New("username/email or password incorrect")
	}
	if !util.VerifyPassword(password, user.PasswordHash) {
		return "", errors.New("username/email or password incorrect")
	}
	jwt, err := authUsecase.generateJWT(user)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (auc *AuthUsecase) generateJWT(user *User) (string, error) {
	// 设置 JWT 声明
	registeredClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(auc.jwtConf.expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    auc.jwtConf.issuer,
		Subject:   user.Username,
		ID:        strconv.FormatInt(user.ID, 10),
	}
	claims := MyCustomClaims{
		UserID:           user.ID,
		RegisteredClaims: registeredClaims,
	}
	// HS256 签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(auc.jwtConf.secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
