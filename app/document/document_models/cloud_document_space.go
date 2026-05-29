package document_models

import "beaver/common/models"

func (CloudDocumentSpace) TableName() string {
	return "cloud_document_spaces"
}

// CloudDocumentSpace 云文档空间（个人盘、团队盘、知识库等）
// 一棵树归属一个 space_id；cloud_documents.space_id 关联本表
type CloudDocumentSpace struct {
	models.Model
	SpaceID     string `gorm:"column:space_id;size:64;uniqueIndex;not null" json:"spaceId"` // 空间业务 ID，如 user:{userId} team:{teamId} wiki:{wikiId}
	SpaceType   int8   `gorm:"column:space_type;not null;default:1" json:"spaceType"`      // 空间类型：1个人 2团队 3知识库
	OwnerID     string `gorm:"column:owner_id;size:64;not null;index" json:"ownerId"`      // 空间所有者 user_id
	CreatorID   string `gorm:"column:creator_id;size:64;not null" json:"creatorId"`        // 空间创建人 user_id
	Name        string `gorm:"column:name;size:128;not null" json:"name"`                  // 空间名称，如「我的云文档」「产品组」
	Description string `gorm:"column:description;size:512" json:"description,omitempty"`   // 空间描述
	Icon        string `gorm:"column:icon;size:512" json:"icon,omitempty"`                 // 空间图标 URL
	DefaultPerm int    `gorm:"column:default_perm;not null;default:1" json:"defaultPerm"`  // 空间成员默认文档权限：1阅读 2编辑 3管理
	Status      int8   `gorm:"column:status;not null;default:1" json:"status"`           // 空间状态：1正常 2删除
}
