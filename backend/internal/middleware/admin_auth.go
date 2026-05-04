package middleware

import (
	"nosocial/internal/bootstrap"
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const AdminIDKey = "admin_id"
const AdminRoleKey = "admin_role"
const AdminUsernameKey = "admin_username"

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 优先从 Cookie 读取
		if cookie, err := c.Cookie("admin_token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// 回退到 Authorization Header
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
			tokenString = parts[1]
		}

		claims, err := jwt.ParseAdminToken(tokenString)
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				response.Unauthorized(c, "登录已过期")
			case jwt.ErrTokenRevoked:
				response.Unauthorized(c, "登录已被吊销，请重新登录")
			default:
				response.Unauthorized(c, "无效的认证信息")
			}
			c.Abort()
			return
		}

		// Casbin权限检查（fail-closed：Enforcer 不可用时拒绝所有请求）
		if bootstrap.Enforcer == nil {
			response.ServerError(c, "权限服务不可用")
			c.Abort()
			return
		}
		obj := c.Request.URL.Path
		act := c.Request.Method
		role := getRoleName(claims.Role)
		ok, _ := bootstrap.Enforcer.Enforce(role, obj, act)
		if !ok {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Set(AdminIDKey, claims.AdminID)
		c.Set(AdminRoleKey, claims.Role)
		c.Set(AdminUsernameKey, claims.Username)
		c.Next()
	}
}

func getRoleName(role int8) string {
	switch role {
	case 3:
		return "superadmin"
	case 1:
		return "manager"
	case 2:
		return "staff"
	default:
		return "" // 未知角色，Casbin 将拒绝所有权限
	}
}

// GetAdminID 从上下文获取管理员ID
func GetAdminID(c *gin.Context) uint32 {
	adminID, _ := c.Get(AdminIDKey)
	if id, ok := adminID.(uint32); ok {
		return id
	}
	return 0
}
