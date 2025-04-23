package utils

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/internal/g"
	"net"
)

const (
	AccessTokenName = "Authorization"
	ClaimsName      = "claims"
)

func ClearAccessToken(c *gin.Context) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie(AccessTokenName, "", -1, "/", "", false, false)
	} else {
		c.SetCookie(AccessTokenName, "", -1, "/", host, false, false)
	}
}

func SetAccessToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie(AccessTokenName, token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie(AccessTokenName, token, maxAge, "/", host, false, false)
	}
}

func GetAccessToken(c *gin.Context) string {
	token := c.Request.Header.Get(AccessTokenName)
	if token == "" {
		j := NewJWT(g.Config.Jwt.Secret)
		token, _ = c.Cookie(AccessTokenName)
		claims, err := j.ParseToken(token)
		if err != nil {
			g.Log.Error("重新写入cookie token失败 " + err.Error())
			return token
		}
		_ = claims
	}
	return token
}

func GetAccessClaims(c *gin.Context) (*CustomClaims, error) {
	token := GetAccessToken(c)
	return ParseToken(token)
}

func ParseToken(token string) (*CustomClaims, error) {
	j := NewJWT(g.Config.Jwt.Secret)
	claims, err := j.ParseToken(token)
	if err != nil {
		g.Log.Error("解析token失败, 请检查token是否过期 " + err.Error())
		return nil, err
	}
	return claims, nil
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint64 {
	if claims, exists := c.Get(ClaimsName); !exists {
		if cl, err := GetAccessClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) int32 {
	if claims, exists := c.Get(ClaimsName); !exists {
		if cl, err := GetAccessClaims(c); err != nil {
			return 0
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.AuthorityId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *CustomClaims {
	if claims, exists := c.Get(ClaimsName); !exists {
		if cl, err := GetAccessClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse
	}
}

// GetUserName 从Gin的Context中获取从jwt解析出来的用户名
func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get(ClaimsName); !exists {
		if cl, err := GetAccessClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.Username
	}
}

func CreateAccessToken(bc BaseClaims) (token string, claims CustomClaims, err error) {
	j := NewJWT(g.Config.Jwt.Secret)
	claims = j.CreateClaims(bc, g.Config.Jwt.Expire.AsDuration(), TokenTypeAccess)
	token, err = j.CreateToken(claims)
	return
}

func CreateRefreshToken(bc BaseClaims) (token string, claims CustomClaims, err error) {
	j := NewJWT(g.Config.Jwt.Secret)
	claims = j.CreateClaims(bc, g.Config.Jwt.RefreshExpire.AsDuration(), TokenTypeRefresh)
	token, err = j.CreateToken(claims)
	return
}
