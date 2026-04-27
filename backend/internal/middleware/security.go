package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders 添加安全 HTTP 响应头
func SecurityHeaders() gin.HandlerFunc {
	// 判断是否为开发环境
	isDev := os.Getenv("NOSOCIAL_APP_MODE") != "release"

	return func(c *gin.Context) {
		// 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		// 启用 XSS 过滤器
		c.Header("X-XSS-Protection", "1; mode=block")
		// 强制 HTTPS（仅生产环境生效）
		if !isDev {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self';")
		// 禁止 referrer 泄露到外部
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// 控制权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
