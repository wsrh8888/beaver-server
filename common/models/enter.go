package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Model struct {
	Id        uint       `gorm:"primaryKey;autoIncrement" json:"id"` // 自增主键
	CreatedAt CustomTime `json:"createdAt"`
	UpdatedAt CustomTime `json:"updatedAt"`
}

type CustomTime time.Time

const layout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	// 先尝试RFC3339格式（兼容带T和时区的数据）
	t, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		*ct = CustomTime(t)
		return nil
	}

	// 可选：保留对旧格式的兼容
	t, err = time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		*ct = CustomTime(t)
		return nil
	}

	return fmt.Errorf("invalid time format: %s", s)
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ct).Format(layout) + `"`), nil
}

func (ct CustomTime) String() string {
	return time.Time(ct).Format(layout)
}

// 添加这两个方法到 CustomTime 类型
func (ct CustomTime) Value() (driver.Value, error) {
	return time.Time(ct), nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(time.Time); ok {
		*ct = CustomTime(v)
		return nil
	}
	return fmt.Errorf("无法扫描 %T 到 CustomTime", value)
}

type PageInfo struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Key   string `json:"key"`
}
