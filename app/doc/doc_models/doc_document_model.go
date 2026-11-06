package doc_models

import (
	"time"

	"beaver/common/models"
)

// DocDocument 在线文档（权威正文源）。表 doc_documents。
// 权威边界：content_json 为 block 结构化快照（渲染/查询用），content_text 为派生纯文本（喂 RAG 投影），
// content_digest 判重避免重复投影索引。version 单调递增，每次内容保存生成一条 doc_revisions。
// knowledge_doc_id 记录该文档在 knowledge 域对应的 Document（RAG 投影行），便于保存/删除时同步。
type DocDocument struct {
	models.Model
	DocID           string `gorm:"column:doc_id;type:varchar(64);uniqueIndex;not null;comment:文档业务ID" json:"docId"`
	NodeID          string `gorm:"column:node_id;type:varchar(64);not null;default:'';index;comment:所属树节点" json:"nodeId"`
	SpaceID         string `gorm:"column:space_id;type:varchar(64);not null;index;comment:所属空间" json:"spaceId"`
	Title           string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	ContentJSON     string `gorm:"column:content_json;type:longtext;comment:block 结构化快照 JSON" json:"contentJson"`
	ContentText     string `gorm:"column:content_text;type:longtext;comment:派生纯文本（供 RAG）" json:"contentText"`
	ContentDigest   string `gorm:"column:content_digest;type:varchar(64);not null;default:'';comment:正文哈希，判重/免重复投影" json:"contentDigest"`
	Version         int    `gorm:"type:int;not null;default:1;comment:版本号，内容每次保存 +1" json:"version"`
	OwnerUserID     string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户" json:"ownerUserId"`
	LastEditedBy    string `gorm:"column:last_edited_by;type:varchar(64);not null;default:'';comment:最后编辑者" json:"lastEditedBy"`
	RagEnabled      bool   `gorm:"column:rag_enabled;type:tinyint;not null;default:1;comment:是否投影进 RAG" json:"ragEnabled"`
	KnowledgeDocID  string `gorm:"column:knowledge_doc_id;type:varchar(64);not null;default:'';comment:对应 RAG 投影 knowledge document" json:"knowledgeDocId"`
	Status          int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间" json:"deletedAt"`
}

func (DocDocument) TableName() string { return "doc_documents" }

func (d *DocDocument) IsDeleted() bool {
	return !d.DeletedAt.IsZero() && d.DeletedAt.Year() > 1970
}
