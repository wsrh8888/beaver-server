package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type SystemMsg struct {
	Type int `json:"type"` //消息类型 1:涉黄 2:涉恐 3:涉政 4:不正确的消息
}

/**
 * @description: 取出来的时候的数据
 */
func (c *SystemMsg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

/**
 * @description: 入库的数据
 */
func (c *SystemMsg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
