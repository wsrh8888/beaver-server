package doc_models

import (
	"time"

	"beaver/common/models"
)

// DocNode 文档树节点（page|folder）。表 doc_nodes。
// 空间内的层级树：parent_id 空为空间根。type=page 时 doc_id 关联一篇 doc_documents；folder 时 doc_id 为空。
type DocNode struct {
	models.Model
	NodeID    string `gorm:"column:node_id;type:varchar(64);uniqueIndex;not null;comment:节点业务ID" json:"nodeId"`
	SpaceID   string `gorm:"column:space_id;type:varchar(64);not null;index;comment:所属空间" json:"spaceId"`
	ParentID  string `gorm:"column:parent_id;type:varchar(64);not null;default:'';index;comment:父节点，空为空间根" json:"parentId"`
	Title     string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Icon      string `gorm:"type:varchar(64);not null;default:'';comment:图标" json:"icon"`
	Type      string `gorm:"type:varchar(16);not null;default:'page';comment:page|folder" json:"type"`
	DocID     string `gorm:"column:doc_id;type:varchar(64);not null;default:'';index;comment:page 关联文档ID" json:"docId"`
	Sort      int    `gorm:"type:int;not null;default:0;comment:同级排序" json:"sort"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间" json:"deletedAt"`
}

func (DocNode) TableName() string { return "doc_nodes" }

func (n *DocNode) IsDeleted() bool {
	return !n.DeletedAt.IsZero() && n.DeletedAt.Year() > 1970
}
