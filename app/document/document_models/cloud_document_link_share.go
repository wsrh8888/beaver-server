package document_models

import "beaver/common/models"

func (CloudDocumentLinkShare) TableName() string {
	return "cloud_document_link_shares"
}

// CloudDocumentLinkShare 文档外链分享配置
// 对外 URL 使用 share_token，避免直接暴露 doc_id
type CloudDocumentLinkShare struct {
	models.Model
	DocID      string             `gorm:"column:doc_id;size:64;not null;index" json:"docId"`            // 关联 cloud_documents.doc_id
	ShareToken string             `gorm:"column:share_token;size:64;uniqueIndex;not null" json:"shareToken"` // 外链 token，拼分享 URL 用
	Perm       int                `gorm:"column:perm;not null;default:1" json:"perm"`                   // 外链访问权限：1阅读 2编辑
	Scope      int8               `gorm:"column:scope;not null;default:1" json:"scope"`                 // 可见范围：1组织内 2互联网公开
	Password   string             `gorm:"column:password;size:128" json:"-"`                              // 访问密码（存哈希，不返回给前端）
	ExpireAt   *models.CustomTime `gorm:"column:expire_at" json:"expireAt,omitempty"`                   // 链接过期时间，空表示永久
	CreatedBy  string             `gorm:"column:created_by;size:64;not null" json:"createdBy"`          // 创建该分享链接的用户 user_id
	Status     int8               `gorm:"column:status;not null;default:1" json:"status"`               // 链接状态：1启用 0关闭
}
