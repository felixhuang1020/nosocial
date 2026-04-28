package logx

// Package logx 提供日志脱敏工具，统一敏感字段的打印格式，
// 防止 openid / transaction_id / mch_id / 证书信息 等落盘/暴露。

// MaskMid 商户号脱敏：保留前 4 位
// 示例：1612345678 -> 1612****
func MaskMid(s string) string {
	if len(s) <= 4 {
		if s == "" {
			return ""
		}
		return "****"
	}
	return s[:4] + "****"
}

// MaskTx 微信支付交易号脱敏：保留首尾 4 位
// 示例：4200001234202501012345678901 -> 4200****8901
func MaskTx(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// MaskOpenID openid 脱敏：保留前 4 位与后 4 位
// 示例：o4GgaunHGD4K7XHQh4u9Iy9FlfVc -> o4Gg****lfVc
func MaskOpenID(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// MaskPhone 手机号脱敏：中间 4 位星号
func MaskPhone(s string) string {
	if len(s) < 7 {
		return "***"
	}
	return s[:3] + "****" + s[len(s)-4:]
}

// MaskID 身份证号/银行卡号脱敏：保留首尾 4 位
func MaskID(s string) string {
	if len(s) < 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// MaskQuery 对 URL RawQuery 中的敏感参数进行脱敏，非敏感字段原样返回。
// 敏感关键词不区分大小写：token/authorization/password/secret/code/signature/key/openid/phone/mobile。
func MaskQuery(raw string) string {
	if raw == "" {
		return ""
	}
	var b []byte
	first := true
	start := 0
	for i := 0; i <= len(raw); i++ {
		if i == len(raw) || raw[i] == '&' {
			pair := raw[start:i]
			eq := indexByte(pair, '=')
			var k, v string
			if eq < 0 {
				k = pair
			} else {
				k = pair[:eq]
				v = pair[eq+1:]
			}
			if isSensitiveKey(k) && v != "" {
				v = "***"
			}
			if !first {
				b = append(b, '&')
			}
			first = false
			b = append(b, k...)
			if eq >= 0 {
				b = append(b, '=')
				b = append(b, v...)
			}
			start = i + 1
		}
	}
	return string(b)
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

var sensitiveQueryKeys = map[string]struct{}{
	"token": {}, "authorization": {}, "password": {}, "pwd": {},
	"secret": {}, "code": {}, "signature": {}, "sign": {}, "key": {},
	"openid": {}, "phone": {}, "mobile": {}, "access_token": {},
	"ossaccesskeyid": {}, "policy": {},
}

func isSensitiveKey(k string) bool {
	low := make([]byte, len(k))
	for i := 0; i < len(k); i++ {
		c := k[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		low[i] = c
	}
	_, ok := sensitiveQueryKeys[string(low)]
	return ok
}
