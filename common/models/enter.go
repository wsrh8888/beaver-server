package models

import "time"

type Model struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 自增主键
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PageInfo struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
	Key   string `json:"key"`
}
