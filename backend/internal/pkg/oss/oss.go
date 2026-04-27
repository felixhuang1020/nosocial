package oss

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"nosocial/config"
	"strings"
	"time"
)

// PostPolicy OSS 服务端签名直传响应（不暴露永久密钥）
type PostPolicy struct {
	Host           string `json:"host"`           // 上传地址
	Key            string `json:"key"`            // 文件路径
	Policy         string `json:"policy"`         // Base64 编码的 Policy
	OSSAccessKeyID string `json:"OSSAccessKeyId"` // 仅 AccessKeyID
	Signature      string `json:"signature"`      // 签名
	Expires        string `json:"expires"`        // 过期时间
	Endpoint       string `json:"endpoint"`       // OSS Endpoint
	Bucket         string `json:"bucket"`         // Bucket 名称
}

// SetObjectInline 通过 CopyObject 将对象的 Content-Disposition 设为 inline。
// 由于 OSS PostObject 不支持通过 form 字段设置 Content-Disposition，
// 需要在上传完成后单独调用此函数修复。
func SetObjectInline(key string) error {
	cfg := config.C.OSS
	bucket := cfg.Bucket
	endpoint := cfg.Endpoint
	date := time.Now().UTC().Format(http.TimeFormat)

	// 只有 x-oss-* 头参与 CanonicalizedOSSHeaders
	signHeaders := []struct{ K, V string }{
		{"x-oss-copy-source", fmt.Sprintf("/%s/%s", bucket, key)},
		{"x-oss-metadata-directive", "REPLACE"},
	}
	var canonicalized string
	for _, h := range signHeaders {
		canonicalized += fmt.Sprintf("%s:%s\n", h.K, h.V)
	}
	canonicalResource := fmt.Sprintf("/%s/%s", bucket, key)

	stringToSign := fmt.Sprintf("PUT\n\napplication/x-www-form-urlencoded\n%s\n%s%s", date, canonicalized, canonicalResource)
	h := hmac.New(sha1.New, []byte(cfg.AccessKeySecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	url := fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, key)
	req, err := http.NewRequest(http.MethodPut, url, io.NopCloser(strings.NewReader("")))
	if err != nil {
		return err
	}
	req.Header.Set("Host", fmt.Sprintf("%s.%s", bucket, endpoint))
	req.Header.Set("Date", date)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-oss-copy-source", fmt.Sprintf("/%s/%s", bucket, key))
	req.Header.Set("x-oss-metadata-directive", "REPLACE")
	req.Header.Set("Content-Disposition", "inline")
	req.Header.Set("Authorization", fmt.Sprintf("OSS %s:%s", cfg.AccessKeyID, signature))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("CopyObject failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

// GeneratePostSignature 生成OSS PostObject签名（服务端签名直传，不泄露永久密钥）
// ext 为可选扩展名（不带点），以便对象 URL 和 OSS Browser 正确识别图片类型
func GeneratePostSignature(dir, ext string) (*PostPolicy, error) {
	cfg := config.C.OSS

	// 生成唯一文件名（可选带扩展名）
	key := fmt.Sprintf("%s/%d_%s", dir, time.Now().Unix(), generateRandomString(8))
	if ext != "" {
		key = key + "." + ext
	}

	// 设置过期时间为1小时
	expireTime := time.Now().Add(1 * time.Hour)

	// 构建Policy
	policyMap := fmt.Sprintf(
		`{"expiration":"%s","conditions":[["content-length-range",0,104857600],["starts-with","$key","%s/"]]}`,
		expireTime.UTC().Format("2006-01-02T15:04:05Z"),
		dir,
	)
	policy := base64.StdEncoding.EncodeToString([]byte(policyMap))

	// 计算Signature
	h := hmac.New(sha1.New, []byte(cfg.AccessKeySecret))
	h.Write([]byte(policy))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return &PostPolicy{
		Host:           fmt.Sprintf("https://%s.%s", cfg.Bucket, cfg.Endpoint),
		Key:            key,
		Policy:         policy,
		OSSAccessKeyID: cfg.AccessKeyID,
		Signature:      signature,
		Expires:        expireTime.Format("2006-01-02T15:04:05Z"),
		Endpoint:       cfg.Endpoint,
		Bucket:         cfg.Bucket,
	}, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 回退到时间种子（理论上不会发生）
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
