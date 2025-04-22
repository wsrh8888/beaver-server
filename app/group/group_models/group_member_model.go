package group_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
)

type GroupMemberModel struct {
	models.Model
	GroupID         string                `gorm:"size:64" json:"groupId"`                     // 群Id
	UserID          string                `json:"userId"`                                     // 用户Id
	MemberNickname  string                `gorm:"size:32" json:"memberNickname"`              // 群成员昵称
	Role            int8                  `json:"role"`                                       // 角色 1:群主 2、管理员 3、普通成员
	ProhibitionTime *int                  `json:"prohibitionTime"`                            // 禁言时间 单位分钟
	UserModel       user_models.UserModel `gorm:"foreignKey:UserID;references:UUID" json:"-"` // 用户信息
	InviterID       string                `gorm:"size:64" json:"inviterId"`                   // 邀请人ID
	Status          int8                  `gorm:"default:1" json:"status"`                    // 成员状态：1正常 2退出 3被踢出
	NotifyLevel     int8                  `gorm:"default:1" json:"notifyLevel"`               // 消息通知级别：1接收所有 2接收@消息 3不接收
	DisplayName     string                `gorm:"size:32" json:"displayName"`                 // 群内显示名称
}
