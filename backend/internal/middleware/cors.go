package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// allowedOrigins 生产环境允许的域名列表
// 开发环境可通过环境变量 NOSOCIAL_CORS_ORIGINS 覆盖，多个用逗号分隔
var allowedOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
	"http://localhost:5173",
	"http://127.0.0.1:5173",
}

func CORS() gin.HandlerFunc {
	origins := allowedOrigins
	if envOrigins := os.Getenv("NOSOCIAL_CORS_ORIGINS"); envOrigins != "" {
		origins = strings.Split(envOrigins, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
