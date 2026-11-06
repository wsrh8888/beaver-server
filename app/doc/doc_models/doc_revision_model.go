package doc_models

import "beaver/common/models"

// DocRevision 用户可见版本历史。表 doc_revisions。
// 每次正文保存/回滚生成一条，保存该版本 content_json 快照，支持列表/查看/diff/回滚。
type DocRevision struct {
	models.Model
	RevisionID   string `gorm:"column:revision_id;type:varchar(64);uniqueIndex;not null;comment:版本业务ID" json:"revisionId"`
	DocID        string `gorm:"column:doc_id;type:varchar(64);not null;index;comment:所属文档" json:"docId"`
	Version      int    `gorm:"type:int;not null;comment:对应文档版本号" json:"version"`
	Title        string `gorm:"type:varchar(255);not null;default:'';comment:该版本标题快照" json:"title"`
	Label        string `gorm:"type:varchar(128);not null;default:'';comment:用户命名" json:"label"`
	ContentJSON  string `gorm:"column:content_json;type:longtext;comment:该版本 block 快照" json:"contentJson"`
	CreatedBy    string `gorm:"column:created_by;type:varchar(64);not null;default:'';comment:产生者" json:"createdBy"`
}

func (DocRevision) TableName() string { return "doc_revisions" }
