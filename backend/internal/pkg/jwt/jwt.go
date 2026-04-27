package jwt

import (
	"errors"
	"fmt"
	"nosocial/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

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
		return claims, nil
	}
	return nil, ErrInvalidToken
}
