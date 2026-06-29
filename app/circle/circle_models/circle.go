package circle_models

import "beaver/common/models"

// CircleModel 圈子主表
type CircleModel struct {
	models.Model
	CircleID    string `gorm:"column:circle_id;size:64;uniqueIndex;not null" json:"circleId"` // 圈子唯一ID
	Name        string `gorm:"size:64;not null" json:"name"`                                  // 圈子名称
	Description string `gorm:"type:text" json:"description"`                                  // 圈子简介
	Avatar      string `gorm:"size:256" json:"avatar"`                                        // 圈子头像
	CreatorID   string `gorm:"size:64;not null;index" json:"creatorId"`                       // 创建者用户ID
	JoinType    int8   `gorm:"not null;default:0" json:"joinType"`                            // 加入方式：0=自由加入 1=审批加入
	MemberCount int64  `gorm:"not null;default:0" json:"memberCount"`                         // 成员数量（冗余字段）
	PostCount   int64  `gorm:"not null;default:0" json:"postCount"`                           // 帖子数量（冗余字段）
	Version     int64  `gorm:"not null;default:0;index" json:"version"`                       // 版本号，用于客户端增量同步
	IsDeleted   bool   `gorm:"not null;default:false;index" json:"isDeleted"`                 // 软删除标记
}
