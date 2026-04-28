package jwt

import (
	"context"
	"errors"
	"fmt"
	"nosocial/config"
	"nosocial/internal/bootstrap"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenRevoked = errors.New("token revoked")
)

const revokePrefix = "nosocial:jwt:revoked:"

// Revoke 将给定 jti 加入 Redis 黑名单（ttl 建议与 Token 剩余有效期保持一致）
// Redis 不可用时入 fail-open：返回错误但不影响已颁发 token 的验证流程。
func Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" || bootstrap.Redis == nil {
		return nil
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return bootstrap.Redis.Set(ctx, revokePrefix+jti, "1", ttl).Err()
}

// IsRevoked 检查 jti 是否已被吊销。Redis 不可用时返回 false（fail-open）避免全站登录故障。
func IsRevoked(ctx context.Context, jti string) bool {
	if jti == "" || bootstrap.Redis == nil {
		return false
	}
	cnt, err := bootstrap.Redis.Exists(ctx, revokePrefix+jti).Result()
	if err != nil {
		return false
	}
	return cnt > 0
}

func newJTI() string {
	return uuid.NewString()
}

// WXClaims 小程序用户JWT Claims
type WXClaims struct {
	UserID uint64 `json:"user_id"`
	Openid string `json:"openid"`
	jwt.RegisteredClaims
}

// AdminClaims 管理员JWT Claims
type AdminClaims struct {
	AdminID  uint32 `json:"admin_id"`
	Username string `json:"username"`
	Role     int8   `json:"role"`
	jwt.RegisteredClaims
}

// GenerateWXToken 生成小程序用户Token
func GenerateWXToken(userID uint64, openid string) (string, error) {
	cfg := config.C.JWT
	claims := WXClaims{
		UserID: userID,
		Openid: openid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        newJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.WXExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.WXSecret))
}

// ParseWXToken 解析小程序用户Token
func ParseWXToken(tokenString string) (*WXClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WXClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 强制验证签名算法，防止算法切换攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.C.JWT.WXSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(*WXClaims); ok && token.Valid {
		if IsRevoked(context.Background(), claims.ID) {
			return nil, ErrTokenRevoked
		}
		return claims, nil
	}
	return nil, ErrInvalidToken
}

// GenerateAdminToken 生成管理员Token
func GenerateAdminToken(adminID uint32, username string, role int8) (string, error) {
	cfg := config.C.JWT
	claims := AdminClaims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        newJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.AdminExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.AdminSecret))
}

// ParseAdminToken 解析管理员Token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 强制验证签名算法，防止算法切换攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.C.JWT.AdminSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		if IsRevoked(context.Background(), claims.ID) {
			return nil, ErrTokenRevoked
		}
		return claims, nil
	}
	return nil, ErrInvalidToken
}
