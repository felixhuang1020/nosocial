package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateInviteCode 生成邀请码（6位字母数字）
func GenerateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

// GenerateOrderNo 生成订单编号
func GenerateOrderNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102"), randomInt(100000, 999999))
}

// GenerateCouponNo 生成优惠券编号
func GenerateCouponNo() string {
	return fmt.Sprintf("CP%s%06d", time.Now().Format("20060102"), randomInt(100000, 999999))
}

// HashPassword bcrypt加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Today 返回今天日期字符串
func Today() string {
	return time.Now().Format("2006-01-02")
}

// TodayMonthDay 返回今天月-日字符串 (MM-DD)
func TodayMonthDay() string {
	return time.Now().Format("01-02")
}

// ParseDate 解析日期字符串
func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func randomInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}
