package doc_models

import (
	"time"

	"beaver/common/models"
)

// DocSpace 文档空间（在线文档顶层容器）。表 doc_spaces。
// 单租户：无 org_id / product_scope。owner_user_id 个人空间归属；team/im 类型可空由成员权限决定可见性。
type DocSpace struct {
	models.Model
	SpaceID     string `gorm:"column:space_id;type:varchar(64);uniqueIndex;not null;comment:空间业务ID" json:"spaceId"`
	Name        string `gorm:"type:varchar(128);not null;comment:空间名" json:"name"`
	Kind        string `gorm:"type:varchar(32);not null;default:'owner';comment:owner个人 team团队 im会话内嵌" json:"kind"`
	OwnerUserID string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户" json:"ownerUserId"`
	ParentID    string `gorm:"column:parent_id;type:varchar(64);not null;default:'';index;comment:父空间ID，空为根" json:"parentId"`
	Icon        string `gorm:"type:varchar(64);not null;default:'';comment:图标标识" json:"icon"`
	Description string `gorm:"type:varchar(512);not null;default:'';comment:描述" json:"description"`
	Status      int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
	Sort        int    `gorm:"type:int;not null;default:0;comment:排序" json:"sort"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间，1970 为未删" json:"deletedAt"`
}

func (DocSpace) TableName() string { return "doc_spaces" }

// IsDeleted 软删判定：deleted_at 非零值即视为已删。
func (s *DocSpace) IsDeleted() bool {
	return !s.DeletedAt.IsZero() && s.DeletedAt.Year() > 1970
}
