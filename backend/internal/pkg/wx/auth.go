// Package wx 封装微信小程序基础接口
package wx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// isPlaceholderWxCred 判断配置值是否为示例占位符（例如 your_wx_appid / __REPLACE__）
func isPlaceholderWxCred(s string) bool {
	if s == "" {
		return true
	}
	lower := strings.ToLower(strings.TrimSpace(s))
	if strings.HasPrefix(lower, "your_") || strings.Contains(lower, "replace") || strings.Contains(lower, "placeholder") {
		return true
	}
	return false
}

// 可通过 SetCode2SessionURL 替换（测试时使用）
var code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

// DefaultHTTPClient 默认 HTTP 客户端（可被替换用于测试）
var DefaultHTTPClient = &http.Client{Timeout: 8 * time.Second}

// Code2SessionResp jscode2session 应答
type Code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// JSCode2Session 使用 code 换取 openid/session_key
// 对外统一返回 "invalid wx code"（避免把微信原始错误外泄），
// 但在服务端日志中会记录真实失败原因（errcode / errmsg / http status），方便排查。
func JSCode2Session(ctx context.Context, appID, secret, code string) (*Code2SessionResp, error) {
	if isPlaceholderWxCred(appID) || isPlaceholderWxCred(secret) {
		log.Printf("[wx.code2session] appid/secret 未正确配置（仍为占位符），请在 backend/.env 中设置 NOSOCIAL_WX_APPID / NOSOCIAL_WX_SECRET")
		return nil, fmt.Errorf("wx appid/secret not configured")
	}
	if strings.TrimSpace(code) == "" {
		log.Printf("[wx.code2session] 前端传入的 code 为空")
		return nil, fmt.Errorf("invalid wx code")
	}

	q := url.Values{}
	q.Set("appid", appID)
	q.Set("secret", secret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")
	reqURL := code2SessionURL + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("[wx.code2session] build request failed: %v", err)
		return nil, fmt.Errorf("invalid wx code")
	}
	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		log.Printf("[wx.code2session] http call failed: %v", err)
		return nil, fmt.Errorf("invalid wx code")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		log.Printf("[wx.code2session] read body failed: %v", err)
		return nil, fmt.Errorf("invalid wx code")
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("[wx.code2session] non-200: status=%d body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("invalid wx code")
	}

	var out Code2SessionResp
	if err := json.Unmarshal(body, &out); err != nil {
		log.Printf("[wx.code2session] unmarshal failed: %v, raw=%s", err, string(body))
		return nil, fmt.Errorf("invalid wx code")
	}
	if out.ErrCode != 0 || out.OpenID == "" {
		// 常见 errcode：40013 invalid appid；40029 invalid code；40125 invalid appsecret；
		// 45011 api minute-quota reach limit；-1 system error
		log.Printf("[wx.code2session] wx api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
		return nil, fmt.Errorf("invalid wx code")
	}
	return &out, nil
}
