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

// AuthEntry is a entry of searching user.
type AuthEntry interface {
	FindUserByNameOrEmail(context.Context, string) (*User, error)
}

type JwtConfig struct {
	secret     string
	expiration time.Duration
	issuer     string
}

// AuthUsecase is a Greeter usecase.
type AuthUsecase struct {
	entry   AuthEntry
	log     *log.Helper
	jwtConf JwtConfig
}

func NewAuthUsecase(entry AuthEntry, logger log.Logger, con *conf.Jwt) *AuthUsecase {
	var exp int32 = 24 // default 24 hours
	if con.Expiration != 0 {
		exp = con.Expiration
	}
	var issuer string = "cascade-oj"
	if con.Issuer != "" {
		issuer = con.Issuer
	}
	return &AuthUsecase{
		entry: entry,
		log:   log.NewHelper(logger),
		jwtConf: JwtConfig{
			secret:     con.Secret,
			expiration: time.Duration(exp) * time.Hour,
			issuer:     issuer,
		},
	}
}

func (auc *AuthUsecase) Signup(ctx context.Context, username string, email string, password string) error {
	userWithSameName, err := auc.entry.FindUserByNameOrEmail(ctx, username)
	if err == nil && userWithSameName != nil {
		return errors.New("username already exists")
	}
	userWithSameEmail, err := auc.entry.FindUserByNameOrEmail(ctx, email)
	if err == nil && userWithSameEmail != nil {
		return errors.New("email already registered")
	}
	return nil
}

func (auc *AuthUsecase) Login(ctx context.Context, usernameOrEmail string, password string) (string, error) {
	user, err := auc.entry.FindUserByNameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return "", errors.New("username/email or password incorrect")
	}
	if !util.VerifyPassword(password, user.PasswordHash) {
		return "", errors.New("username/email or password incorrect")
	}
	jwt, err := auc.generateJWT(user)
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
