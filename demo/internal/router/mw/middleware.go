package mw

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/common/resp"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/g"
	"github/invokerw/gintos/demo/internal/pkg/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetAccessToken(c)
		if token == "" {
			resp.NoAuth(c, "未登录或非法访问")
			c.Abort()
			return
		}

		j := utils.NewJWT(g.Config.Jwt.Secret)
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				resp.NoAuth(c, "授权已过期")
				utils.ClearAccessToken(c)
				c.Abort()
				return
			}
			resp.NoAuth(c, err.Error())
			utils.ClearAccessToken(c)
			c.Abort()
			return
		}

		c.Set(utils.ClaimsName, claims)
		c.Next()
	}
}

func Authorize(limit common.UserAuthority) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := utils.GetUserInfo(c)
		if userClaims == nil {
			resp.NoAuth(c, "未登录或非法访问")
			c.Abort()
			return
		}

		if userClaims.AuthorityId > int32(limit) {
			resp.NoAuth(c, "没有权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
