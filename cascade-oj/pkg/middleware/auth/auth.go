package auth

import (
	"context"
	"encoding/json"
	"strings"

	"cascade-oj/ent/user"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("UNAUTHORIZED", "JWT token is missing")
	ErrUnSupportSigningMethod = errors.Unauthorized("UNAUTHORIZED", "Wrong signing method")
	ErrSignToken              = errors.Unauthorized("UNAUTHORIZED", "Can not sign token. Is the key correct?")
	ErrTokenExpired           = errors.Unauthorized("UNAUTHORIZED", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("UNAUTHORIZED", "Fail to parse JWT token ")
	ErrTokenInvalid           = errors.Unauthorized("UNAUTHORIZED", "Token is invalid")
	ErrMissingClaims          = errors.Unauthorized("UNAUTHORIZED", "Claims is missing")
	ErrWrongContext           = errors.Unauthorized("UNAUTHORIZED", "Wrong context for middleware")
)

type Claims struct {
	UserID int64     `json:"user_id"`
	Role   user.Role `json:"role"`
}

// Authorization
const (
	authorizationKey string = "token"
)

// Roles
const (
	RoleCompetitor string = "competitor"
	RoleCreator    string = "creator"
	RoleAdmin      string = "admin"
)

// Auth middleware
//
// params:
//   - secret: JWT secret key
//   - defaultRule: default access rule for unspecified APIs
//   - customApiMap: map of API paths to access rules
func Auth(secret string, defaultRule string, customApiMap map[string]string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			token, err := getTokenFromTrans(ctx)
			if err != nil {
				return nil, err
			}
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, ErrUnSupportSigningMethod
				}
				return []byte(secret), nil
			}
			claims, err := parseToken(token, keyFunc)
			if err != nil {
				return nil, err
			}
			claimsInfo, err := parseClaims(*claims)
			if err != nil {
				return nil, err
			}
			log.Infof("received request from user ID: %d", claimsInfo.UserID)

			//put claims into context so that other service could retrieve it
			ctx = context.WithValue(ctx, "userInfo", claimsInfo)

			// check the role of user
			trans, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, ErrWrongContext
			}
			operation := trans.Operation()
			role := defaultRule
			if customApiMap != nil && customApiMap[operation] != "" {
				role = customApiMap[operation]
			}
			switch role {
			case RoleCompetitor, RoleCreator, RoleAdmin:
				if !strings.HasPrefix(claimsInfo.Role.String(), role) {
					return nil, errors.Unauthorized("PERMISSION_DENIED", "Wrong role")
				}
			default:
				return nil, errors.InternalServer("INTERNAL_SERVER", "Wrong rule")
			}

			return handler(ctx, req)
		}
	}
}

// get token from context header
func getTokenFromTrans(ctx context.Context) (string, error) {
	header, ok := transport.FromServerContext(ctx)
	if ok {
		return header.RequestHeader().Get(authorizationKey), nil
	}
	log.Warn("Can not get token from transport")
	return "", ErrMissingJwtToken
}

func parseToken(token string, keyFunc jwt.Keyfunc) (*jwt.Claims, error) {
	var (
		tokenInfo *jwt.Token
		err       error
	)
	tokenInfo, err = jwt.Parse(token, keyFunc, jwt.WithJSONNumber())
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenUnverifiable) {
			return nil, ErrSignToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenParseFail
	}
	if !tokenInfo.Valid {
		return nil, ErrTokenInvalid
	}
	return &tokenInfo.Claims, nil
}

func parseClaims(claims jwt.Claims) (*Claims, error) {
	userID, err := claims.(jwt.MapClaims)["user_id"].(json.Number).Int64()
	if err != nil {
		return nil, ErrMissingClaims
	}
	role, ok := claims.(jwt.MapClaims)["role"].(string)
	if !ok {
		return nil, ErrMissingClaims
	}
	return &Claims{
		UserID: userID,
		Role:   user.Role(role),
	}, nil
}
