package document_models

import "beaver/common/models"

func (CloudDocumentRevision) TableName() string {
	return "cloud_document_revisions"
}

// CloudDocumentRevision 文档历史版本快照（仅 doc_type>=1 的文档类资源）
// 每次保存正文前，将上一版 file_key 写入此表
type CloudDocumentRevision struct {
	models.Model
	DocID      string `gorm:"column:doc_id;size:64;not null;uniqueIndex:idx_doc_rev,priority:1" json:"docId"` // 关联 cloud_documents.doc_id
	Revision   int64  `gorm:"column:revision;not null;uniqueIndex:idx_doc_rev,priority:2" json:"revision"`    // 历史版本号（保存前的版本）
	EditorID   string `gorm:"column:editor_id;size:64;not null" json:"editorId"`                              // 该版本的编辑者 user_id
	ChangeNote string `gorm:"column:change_note;size:256" json:"changeNote,omitempty"`                        // 版本备注，如「自动保存」「发布前快照」
	FileKey    string `gorm:"column:file_key;size:64;not null" json:"fileKey"`                                // 该版本正文 file_key
}
