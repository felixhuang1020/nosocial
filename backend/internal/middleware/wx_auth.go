package middleware

import (
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const WXUserIDKey = "wx_user_id"
const WXOpenidKey = "wx_openid"

// WXAuth 小程序用户认证中间件
func WXAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseWXToken(parts[1])
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Unauthorized(c, "登录已过期")
			} else {
				response.Unauthorized(c, "无效的认证信息")
			}
			c.Abort()
			return
		}

		c.Set(WXUserIDKey, claims.UserID)
		c.Set(WXOpenidKey, claims.Openid)
		c.Next()
	}
}

// GetWXUserID 从上下文获取用户ID
func GetWXUserID(c *gin.Context) uint64 {
	userID, _ := c.Get(WXUserIDKey)
	if id, ok := userID.(uint64); ok {
		return id
	}
	return 0
}
