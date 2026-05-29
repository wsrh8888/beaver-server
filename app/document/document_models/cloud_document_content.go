package document_models

import "beaver/common/models"

func (CloudDocumentContent) TableName() string {
	return "cloud_document_contents"
}

// CloudDocumentContent 文档当前正文指针（与 cloud_documents 1:1）
// 文件夹（doc_type=0）无记录；正文实体在 file 服务，此处只存 file_key
type CloudDocumentContent struct {
	models.Model
	DocID    string `gorm:"column:doc_id;size:64;uniqueIndex;not null" json:"docId"` // 关联 cloud_documents.doc_id
	Revision int64  `gorm:"column:revision;not null;default:1" json:"revision"`      // 当前版本号，与 cloud_documents.revision 保持一致
	FileKey  string `gorm:"column:file_key;size:64;not null" json:"fileKey"`         // 正文文件 key，关联 file_models.FileKey
}
