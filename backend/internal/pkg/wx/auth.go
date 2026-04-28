// Package wx 封装微信小程序基础接口
package wx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
// 失败统一返回 "invalid wx code"（避免把微信原始错误外泄）
func JSCode2Session(ctx context.Context, appID, secret, code string) (*Code2SessionResp, error) {
	if appID == "" || secret == "" {
		return nil, fmt.Errorf("wx appid/secret not configured")
	}
	if strings.TrimSpace(code) == "" {
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
		return nil, fmt.Errorf("invalid wx code")
	}
	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invalid wx code")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return nil, fmt.Errorf("invalid wx code")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid wx code")
	}

	var out Code2SessionResp
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("invalid wx code")
	}
	if out.ErrCode != 0 || out.OpenID == "" {
		return nil, fmt.Errorf("invalid wx code")
	}
	return &out, nil
}
