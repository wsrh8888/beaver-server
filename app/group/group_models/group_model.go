package group_models

import "beaver/common/models"

type GroupModel struct {
	models.Model
	UUID               string             `gorm:"size:64;unique" json:"uuid"`                  // 设置主键 UserID
	Title              string             `gorm:"size:32" json:"title"`                        // 群名
	Abstract           string             `gorm:"size:128" json:"abstract"`                    // 简介
	Avatar             string             `gorm:"size:256" json:"avatar"`                      // 头像
	Creator            string             `gorm:"type:longtext" json:"creator"`                // 创建者
	IsInvite           bool               `json:"isInvite"`                                    // 是否允许被邀请
	IsTemporarySession bool               `json:"isTemporarySession"`                          // 是否是临时会话
	IsProhibition      bool               `json:"isProhibition"`                               // 是否开启全员禁言
	Size               int                `json:"size"`                                        // 群规模
	MemberList         []GroupMemberModel `gorm:"foreignKey:GroupID;references:UUID" json:"-"` // 群成员列表
}
