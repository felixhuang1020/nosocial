package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray PostgreSQL JSONB 字符串数组字段。
// 在 Go 侧表现为 []string，在 DB 侧以 JSONB 存储。
type StringArray []string

// Value 实现 driver.Valuer：序列化为 JSON 写入 DB。
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan 实现 sql.Scanner：从 DB JSONB 反序列化。
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("StringArray.Scan: unsupported type")
	}
	if len(bytes) == 0 {
		*a = nil
		return nil
	}
	return json.Unmarshal(bytes, a)
}
