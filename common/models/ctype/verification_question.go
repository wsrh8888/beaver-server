package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type VerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}

/**
 * @description: 取出来的时候的数据
 */
func (c *VerificationQuestion) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

/**
 * @description: 入库的数据
 */
func (c *VerificationQuestion) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
