package doc_models

import (
	"beaver/common/models"
)

// DocMemberModel 文档成员/权限白名单 (ACL 控制表)
// 只有当文档需要“突破” SpaceID 的默认权限空间时，才在这里增加记录。
// 对标大厂设计：支持高效的反向查询 (即“查询所有分享给我的文档”)。
type DocMemberModel struct {
	models.Model
	DocID    string `gorm:"size:64;uniqueIndex:idx_doc_member" json:"docId"`    // 文档ID
	MemberID string `gorm:"size:64;uniqueIndex:idx_doc_member" json:"memberID"` // 成员ID (UID 或 GID)

	// 成员类型：1:用户, 2:群组
	MemberType int8 `gorm:"default:1" json:"memberType"`

	// 权限角色：1:查看者(Viewer), 2:编辑者(Editor), 3:管理者(Admin/Full Access)
	Role int8 `gorm:"default:1" json:"role"`
}
