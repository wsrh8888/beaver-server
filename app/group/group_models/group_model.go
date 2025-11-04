package group_models

import (
	"beaver/common/models"
)

// 每个群独立版本
type GroupModel struct {
	models.Model
	GroupID   string `gorm:"size:64;unique;index" json:"groupId"`
	Type      int8   `gorm:"default:1" json:"type"`                                               // 群类型：1正常群 2讨论组 ...
	Title     string `gorm:"size:32;index" json:"title"`                                          // 群名
	Avatar    string `gorm:"size:256;default:a9de5548bef8c10b92428fff61275c72.png" json:"avatar"` // 群头像文件名
	CreatorID string `gorm:"size:64;index" json:"creatorId"`                                      // 创建者ID
	Notice    string `gorm:"type:text" json:"notice"`                                             // 当前公告内容
	JoinType  int8   `gorm:"not null;default:0" json:"joinType"`                                  // 0自由加入 1需审批 2不可加入
	Status    int8   `gorm:"default:1" json:"status"`                                             // 群状态：1正常 2冻结 3解散
	Version   int64  `gorm:"not null;default:0;index" json:"version"`                             // 数据同步版本号
}
