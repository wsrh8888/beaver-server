package group_models

import (
	"beaver/common/models"
	"time"
)

type GroupModel struct {
	models.Model
	UUID           string             `gorm:"size:64;unique;index" json:"uuid"`                                    // 设置主键 UserID
	Type           int8               `gorm:"default:1" json:"type"`                                               // 群类型
	Title          string             `gorm:"size:32;index" json:"title"`                                          // 群名
	Abstract       string             `gorm:"size:128" json:"abstract"`                                            // 简介
	Avatar         string             `gorm:"size:256;default:e7be4283-dc79-4db7-b65c-aa335b90bcfb" json:"avatar"` // 头像
	CreatorID      string             `gorm:"size:64;index" json:"creatorID"`                                      // 创建者ID
	Notice         string             `gorm:"type:text" json:"notice"`                                             // 群公告
	Tags           string             `gorm:"size:256" json:"tags"`                                                // 群标签
	MaxMembers     int                `gorm:"default:500" json:"maxMembers"`                                       // 群最大成员数
	CurrentMembers int                `gorm:"default:0" json:"currentMembers"`                                     // 当前成员数
	Status         int8               `gorm:"default:1" json:"status"`                                             // 群状态
	MemberList     []GroupMemberModel `gorm:"foreignKey:GroupID;references:UUID" json:"-"`                         // 群成员列表
	MuteAll        bool               `gorm:"default:false" json:"muteAll"`                                        // 全员禁言状态
	DissolveTime   *time.Time         `json:"dissolveTime"`                                                        // 群解散时间
	Category       string             `gorm:"size:32" json:"category"`                                             // 群分类
	Settings       GroupSettings      `gorm:"embedded" json:"settings"`                                            // 嵌入群设置
}

type GroupSettings struct {
	JoinAuth         int8 `gorm:"default:1" json:"joinAuth"`            // 加入权限
	MemberInvite     bool `gorm:"default:true" json:"memberInvite"`     // 允许成员邀请
	MemberManage     bool `gorm:"default:false" json:"memberManage"`    // 允许成员管理
	MessageArchive   bool `gorm:"default:true" json:"messageArchive"`   // 允许消息存档
	AllowViewHistory bool `gorm:"default:true" json:"allowViewHistory"` // 允许查看历史消息
}
