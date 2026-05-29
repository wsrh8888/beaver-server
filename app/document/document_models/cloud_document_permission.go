package document_models

import "beaver/common/models"

const (
	DocPermView   = 1 // 阅读
	DocPermEdit   = 2 // 编辑
	DocPermManage = 3 // 管理
)

const (
	DocSubjectUser  = 1 // 用户
	DocSubjectGroup = 2 // 群/会话
	DocSubjectDept  = 3 // 部门
)

func (CloudDocumentPermission) TableName() string {
	return "cloud_document_permissions"
}

// CloudDocumentPermission 资源级协作者权限（文件夹或文档均可配置）
type CloudDocumentPermission struct {
	models.Model
	DocID       string `gorm:"column:doc_id;size:64;not null;uniqueIndex:idx_doc_subject,priority:1" json:"docId"`       // 关联 cloud_documents.doc_id
	SubjectType int8   `gorm:"column:subject_type;not null;uniqueIndex:idx_doc_subject,priority:2" json:"subjectType"` // 授权主体类型：1用户 2群 3部门
	SubjectID   string `gorm:"column:subject_id;size:64;not null;uniqueIndex:idx_doc_subject,priority:3" json:"subjectId"` // 授权主体 ID（user_id / 群 ID / 部门 ID）
	Perm        int    `gorm:"column:perm;not null;default:1" json:"perm"`                                             // 权限级别：1阅读 2编辑 3管理
	Inherit     bool   `gorm:"column:inherit;not null;default:true" json:"inherit"`                                    // 是否向下继承给子文件夹/子文档
	GrantedBy   string `gorm:"column:granted_by;size:64" json:"grantedBy,omitempty"`                                   // 授权人 user_id
}
