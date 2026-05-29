package document_models

import "beaver/common/models"

func (CloudDocument) TableName() string {
	return "cloud_documents"
}

// CloudDocument 云文档资源树节点（文件夹 + 各类云文档）
// 仅存元数据：标题、目录位置、权限关联、状态等；正文 file_key 见 cloud_document_contents
type CloudDocument struct {
	models.Model
	DocID        string             `gorm:"column:doc_id;size:64;uniqueIndex;not null" json:"docId"`                           // 资源业务 ID（UUID），对标 UserID；IM 分享卡片用此字段
	OwnerID      string             `gorm:"column:owner_id;size:64;not null;index" json:"ownerId"`                             // 当前归属者 user_id，可转让；默认拥有管理权限
	CreatorID    string             `gorm:"column:creator_id;size:64;not null" json:"creatorId"`                               // 创建人 user_id，创建后不变
	DocType      int                `gorm:"column:doc_type;not null;default:1" json:"docType"`                                 // 0文件夹 1文档 2表格 3幻灯片 4思维笔记；打开哪个编辑器由此决定
	Title        string             `gorm:"column:title;size:256;not null" json:"title"`                                       // 名称（文件夹名或文档标题）
	Revision     int64              `gorm:"column:revision;not null;default:1" json:"revision"`                                // 当前版本号；保存正文时递增，用于乐观锁
	SpaceID      string             `gorm:"column:space_id;size:64;not null;index:idx_space_parent,priority:1" json:"spaceId"` // 所属空间 ID，关联 cloud_document_spaces.space_id
	ParentID     string             `gorm:"column:parent_id;size:64;index:idx_space_parent,priority:2" json:"parentId"`        // 父节点 doc_id；空串表示挂在空间根目录
	SortOrder    int                `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`                           // 同级目录下的排序权重，越小越靠前
	CoverURL     string             `gorm:"column:cover_url;size:512" json:"coverUrl,omitempty"`                               // 封面图 URL，列表/分享卡片展示用；文件夹通常为空
	LastEditorID string             `gorm:"column:last_editor_id;size:64" json:"lastEditorId,omitempty"`                       // 最近一次保存正文的编辑者 user_id
	Status       int8               `gorm:"column:status;not null;default:1" json:"status"`                                    // 业务状态：1正常 2归档（仍可见，偏冷存储）
	DeletedAt    *models.CustomTime `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`                                // 软删除时间；非空表示在回收站，列表默认不展示
	DeletedBy    string             `gorm:"column:deleted_by;size:64" json:"deletedBy,omitempty"`                              // 执行删除操作的用户 user_id
}
