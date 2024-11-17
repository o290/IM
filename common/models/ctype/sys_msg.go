package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

// 系统提示
type SystemMsg struct {
	Type int8 `json:"type"` //违规类型：1：涉黄2：涉恐3：涉政4：不正当言论
}

func (m *SystemMsg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), m)
}
func (m SystemMsg) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}
