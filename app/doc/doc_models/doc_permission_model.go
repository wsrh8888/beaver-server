package doc_models

import "beaver/common/models"

// DocPermission 资源授权（分享/协作者）。表 doc_permissions。
// 可绑定 space 或 doc；主体可为 user|group|dept。role：owner|manager|editor|commenter|viewer。
// owner 由创建时直接写入，不通过 SetPermission 授予。对 doc 无显式授权时，按所属 space 的授权继承（CheckPermission 计算）。
type DocPermission struct {
	models.Model
	ResourceType  string `gorm:"column:resource_type;type:varchar(16);not null;uniqueIndex:idx_perm,priority:1;comment:space|doc" json:"resourceType"`
	ResourceID    string `gorm:"column:resource_id;type:varchar(64);not null;uniqueIndex:idx_perm,priority:2;index;comment:空间或文档ID" json:"resourceId"`
	PrincipalType string `gorm:"column:principal_type;type:varchar(16);not null;uniqueIndex:idx_perm,priority:3;comment:user|group|dept" json:"principalType"`
	PrincipalID   string `gorm:"column:principal_id;type:varchar(64);not null;uniqueIndex:idx_perm,priority:4;index;comment:主体ID" json:"principalId"`
	Role          string `gorm:"type:varchar(16);not null;default:'viewer';comment:owner|manager|editor|commenter|viewer" json:"role"`
}

func (DocPermission) TableName() string { return "doc_permissions" }
