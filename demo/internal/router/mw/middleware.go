package mw

import (
	"errors"
	"github/invokerw/gintos/common/resp"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/g"
	"github/invokerw/gintos/demo/internal/pkg/utils"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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

func CasbinAuth(limit common.UserAuthority, e *casbin.Enforcer) gin.HandlerFunc {
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
		if userClaims.AuthorityId != int32(common.UserAuthority_SYS_ADMIN) && e != nil {
			sub := userClaims.Role
			obj := c.Request.URL.Path
			act := c.Request.Method

			ok, err := e.Enforce(sub, obj, act)
			if err != nil {
				resp.NoAuth(c, "权限校验失败 "+err.Error())
				c.Abort()
			}
			if !ok {
				resp.NoAuth(c, "没有权限")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
