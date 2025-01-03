package user_models

import "beaver/common/models"

type UserModel struct {
	models.Model
	UserId   string `gorm:"size:64;unique" json:"userId"` // 唯一字段，但不是主键
	NickName string `json:"nickName"`
	Password string `json:"password"`
	Avatar   string `gorm:"default:'https://js.ibaotu.com/images/avatar/%E5%A4%B4%E5%83%8F-17.png'" json:"avatar"`
	Abstract string `gorm:"size:32" json:"abstract"`
	Phone    string `gorm:"size:11" json:"phone"`
	Source   int32  `json:"source"`
}
