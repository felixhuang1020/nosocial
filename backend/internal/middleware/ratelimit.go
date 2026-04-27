package middleware

import (
	"net/http"
	"nosocial/internal/pkg/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// 定义不同接口的限流规则
var (
	// 通用限流：每IP 100次/分钟
	generalRate = limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// 登录限流：每IP 5次/分钟
	loginRate = limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}

	// OSS签名限流：每IP 30次/分钟（管理端批量上传时 10/min 偏严）
	uploadRate = limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  30,
	}

	// store 限流器存储后端，通过 InitRateLimitStore 初始化
	store limiter.Store = memory.NewStore()
)

// InitRateLimitStore 初始化限流器存储后端。
// 传入 Redis 客户端则使用 Redis 存储（支持多进程共享），否则使用内存存储。
func InitRateLimitStore(rdb *redis.Client) {
	if rdb == nil {
		store = memory.NewStore()
		return
	}
	rs, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:          "nosocial:rate:",
		MaxRetry:        3,
		CleanUpInterval: 5 * time.Minute,
	})
	if err != nil {
		store = memory.NewStore()
		return
	}
	store = rs
}

// RateLimit 通用限流中间件
// scope 用于区分不同限流器的计数器命名空间，避免多种限流策略共享同一 IP 计数器
func RateLimit(rate limiter.Rate, scope string) gin.HandlerFunc {
	instance := limiter.New(store, rate)
	return func(c *gin.Context) {
		key := getClientKey(c) + ":" + scope
		ctx, cancel := c.Request.Context(), func() {}
		defer cancel()

		limitCtx, err := instance.Get(ctx, key)
		if err != nil {
			response.ErrorWithStatus(c, http.StatusInternalServerError, "限流服务异常")
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(limitCtx.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(limitCtx.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(limitCtx.Reset, 10))

		if limitCtx.Reached {
			response.ErrorWithStatus(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitLogin 登录接口专用限流
func RateLimitLogin() gin.HandlerFunc {
	return RateLimit(loginRate, "login")
}

// RateLimitUpload 上传签名接口专用限流
func RateLimitUpload() gin.HandlerFunc {
	return RateLimit(uploadRate, "upload")
}

// RateLimitGeneral 通用接口限流
func RateLimitGeneral() gin.HandlerFunc {
	return RateLimit(generalRate, "general")
}

// getClientKey 获取客户端标识（优先 X-Forwarded-For，其次 RemoteAddr）
func getClientKey(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	// 处理 X-Forwarded-For 可能包含多个 IP 的情况
	if idx := strings.Index(ip, ","); idx != -1 {
		ip = strings.TrimSpace(ip[:idx])
	}
	return ip
}
