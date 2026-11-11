package doc_models

import "beaver/common/models"

// DocRecentVisit 最近访问。表 doc_recent_visits。同一用户对同一文档只保留一行，重复访问刷新时间。
type DocRecentVisit struct {
	models.Model
	UserID    string `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:idx_recent_user_doc,priority:1;comment:访问者" json:"userId"`
	DocID     string `gorm:"column:doc_id;type:varchar(64);not null;uniqueIndex:idx_recent_user_doc,priority:2;comment:文档ID" json:"docId"`
	VisitedAt string `gorm:"column:visited_at;type:varchar(32);not null;default:'';index;comment:访问时间 RFC3339" json:"visitedAt"`
}

func (DocRecentVisit) TableName() string { return "doc_recent_visits" }

// DocFavorite 收藏。表 doc_favorites。
type DocFavorite struct {
	models.Model
	UserID string `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:idx_fav_user_doc,priority:1;comment:收藏者" json:"userId"`
	DocID  string `gorm:"column:doc_id;type:varchar(64);not null;uniqueIndex:idx_fav_user_doc,priority:2;index;comment:文档ID" json:"docId"`
}

func (DocFavorite) TableName() string { return "doc_favorites" }

// DocTemplate 文档模板。表 doc_templates。
// scope：builtin 内置（owner 为空）/ team 团队 / owner 个人。content_json 是 ProseMirror 文档快照，
// doc_type 决定实例化时创建哪种文档（doc|sheet|board）。
type DocTemplate struct {
	models.Model
	TemplateID  string `gorm:"column:template_id;type:varchar(64);uniqueIndex;not null;comment:模板业务ID" json:"templateId"`
	Name        string `gorm:"column:name;type:varchar(128);not null;comment:模板名" json:"name"`
	Category    string `gorm:"column:category;type:varchar(64);not null;default:'';index;comment:分类" json:"category"`
	Description string `gorm:"column:description;type:varchar(512);not null;default:'';comment:描述" json:"description"`
	Cover       string `gorm:"column:cover;type:varchar(512);not null;default:'';comment:封面标识/URL" json:"cover"`
	ContentJSON string `gorm:"column:content_json;type:longtext;comment:模板 ProseMirror 快照" json:"contentJson"`
	DocType     string `gorm:"column:doc_type;type:varchar(16);not null;default:'doc';comment:doc|sheet|board" json:"docType"`
	Scope       string `gorm:"column:scope;type:varchar(16);not null;default:'owner';comment:builtin|team|owner" json:"scope"`
	OwnerUserID string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户" json:"ownerUserId"`
	UseCount    int    `gorm:"column:use_count;type:int;not null;default:0;comment:使用次数" json:"useCount"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
}

func (DocTemplate) TableName() string { return "doc_templates" }

// DocTrash 回收站条目。表 doc_trash。
// 软删的 space/node/doc 统一在此登记一行，保留删除前的父子关系以便原路径恢复；
// purge 时物理删除该行并清理对应业务行。
type DocTrash struct {
	models.Model
	TrashID      string `gorm:"column:trash_id;type:varchar(64);uniqueIndex;not null;comment:回收站条目ID" json:"trashId"`
	ResourceType string `gorm:"column:resource_type;type:varchar(16);not null;index:idx_trash_user,priority:1;comment:space|node|doc" json:"resourceType"`
	ResourceID   string `gorm:"column:resource_id;type:varchar(64);not null;index;comment:资源业务ID" json:"resourceId"`
	Title        string `gorm:"column:title;type:varchar(255);not null;default:'';comment:删除时标题快照" json:"title"`
	SpaceID      string `gorm:"column:space_id;type:varchar(64);not null;default:'';comment:所属空间" json:"spaceId"`
	ParentID     string `gorm:"column:parent_id;type:varchar(64);not null;default:'';comment:删除时父节点" json:"parentId"`
	DeletedBy    string `gorm:"column:deleted_by;type:varchar(64);not null;default:'';index:idx_trash_user,priority:2;comment:删除人" json:"deletedBy"`
	DeletedAt    string `gorm:"column:deleted_at;type:varchar(32);not null;default:'';comment:删除时间 RFC3339" json:"deletedAt"`
	// PayloadJSON 保存删除时的子树快照（node_ids / doc_ids），供 RestoreTrash 精确恢复；
	// 单文档删除时为空。不暴露给客户端。
	PayloadJSON string `gorm:"column:payload_json;type:longtext;comment:子树快照 JSON" json:"payloadJson"`
}

func (DocTrash) TableName() string { return "doc_trash" }
