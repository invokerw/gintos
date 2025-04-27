package utils

import (
	"errors"
	"github/invokerw/gintos/demo/internal/g"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType int32

const (
	TokenTypeAccess TokenType = iota + 1
	TokenTypeRefresh
)

type BaseClaims struct {
	ID          uint64
	Username    string
	NickName    string
	AuthorityId int32
	Role        string
}

type CustomClaims struct {
	BaseClaims
	TokenType TokenType
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenValid            = errors.New("未知错误")
	TokenExpired          = errors.New("token已过期")
	TokenNotValidYet      = errors.New("token尚未激活")
	TokenMalformed        = errors.New("这不是一个token")
	TokenSignatureInvalid = errors.New("无效签名")
	TokenInvalid          = errors.New("无法处理此token")
)

func NewJWT(key string) *JWT {
	return &JWT{
		SigningKey: []byte(key),
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims, expires time.Duration, tokenType TokenType) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		TokenType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Gintos"},                  // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),   // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)), // 过期时间
			Issuer:    g.Config.Jwt.Issuer,                         // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenValid
}
