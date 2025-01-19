package models

import "time"

type Model struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 自增主键
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ModelAuth struct {
	ID        uint      `gorm:"autoIncrement" json:"id"` // 保留自增 ID，去掉 primaryKey 标签
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PageInfo struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Key   string `json:"key"`
}
