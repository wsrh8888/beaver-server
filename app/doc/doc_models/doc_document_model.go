package doc_models

import (
	"time"

	"beaver/common/models"
)

// DocDocument 在线文档（权威正文源）。表 doc_documents。
// 权威边界：content_json 为 block 结构化快照（渲染/查询用），content_text 为派生纯文本，
// content_digest 判重避免重复索引。version 单调递增，每次内容保存生成一条 doc_revisions。
type DocDocument struct {
	models.Model
	DocID         string `gorm:"column:doc_id;type:varchar(64);uniqueIndex;not null;comment:文档业务ID" json:"docId"`
	NodeID        string `gorm:"column:node_id;type:varchar(64);not null;default:'';index;comment:所属树节点" json:"nodeId"`
	SpaceID       string `gorm:"column:space_id;type:varchar(64);not null;index;comment:所属空间" json:"spaceId"`
	Title         string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	ContentJSON   string `gorm:"column:content_json;type:longtext;comment:block 结构化快照 JSON" json:"contentJson"`
	ContentText   string `gorm:"column:content_text;type:longtext;comment:派生纯文本" json:"contentText"`
	ContentDigest string `gorm:"column:content_digest;type:varchar(64);not null;default:'';comment:正文哈希，判重/免重复索引" json:"contentDigest"`
	Version       int    `gorm:"type:int;not null;default:1;comment:版本号，内容每次保存 +1" json:"version"`
	// Epoch 是协同房间纪元。版本回滚无法在既有 CRDT 状态上「倒退」，因此每次回滚把 epoch +1，
	// 房间名变为 docId#epoch，旧房间的 CRDT 状态整体作废，所有客户端以新房间重新同步。
	Epoch         int    `gorm:"column:epoch;type:int;not null;default:1;comment:协同房间纪元，回滚时递增" json:"epoch"`
	DocType       string `gorm:"column:doc_type;type:varchar(16);not null;default:'doc';index;comment:doc|sheet|board|slide" json:"docType"`
	Icon          string `gorm:"column:icon;type:varchar(64);not null;default:'';comment:文档图标" json:"icon"`
	Cover         string `gorm:"column:cover;type:varchar(512);not null;default:'';comment:封面标识/URL" json:"cover"`
	RagEnabled    bool   `gorm:"column:rag_enabled;type:tinyint(1);not null;default:0;comment:是否允许进入知识库索引" json:"ragEnabled"`
	OwnerUserID   string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户" json:"ownerUserId"`
	LastEditedBy  string `gorm:"column:last_edited_by;type:varchar(64);not null;default:'';comment:最后编辑者" json:"lastEditedBy"`
	Status        int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间" json:"deletedAt"`
}

func (DocDocument) TableName() string { return "doc_documents" }

func (d *DocDocument) IsDeleted() bool {
	return !d.DeletedAt.IsZero() && d.DeletedAt.Year() > 1970
}
