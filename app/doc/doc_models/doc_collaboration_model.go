package doc_models

import (
	"time"

	"beaver/common/models"
)

// DocComment 文档评论 / 回复。表 doc_comments。
// anchor_id 是被评论内容块的稳定 ID（由编辑器写入 PM node attrs），quote 是被选中文本快照，
// 用于正文漂移后仍能定位评论位置。parent_id 为空表示顶层评论，否则为该评论的回复。
// resolve 只作用于顶层评论，其回复随顶层一起关闭。
type DocComment struct {
	models.Model
	CommentID   string `gorm:"column:comment_id;type:varchar(64);uniqueIndex;not null;comment:评论业务ID" json:"commentId"`
	DocID       string `gorm:"column:doc_id;type:varchar(64);not null;index:idx_comment_doc;comment:所属文档" json:"docId"`
	Epoch       int    `gorm:"column:epoch;type:int;not null;default:1;comment:产生时的文档纪元" json:"epoch"`
	ParentID    string `gorm:"column:parent_id;type:varchar(64);not null;default:'';index;comment:父评论ID，空为顶层" json:"parentId"`
	AnchorID    string `gorm:"column:anchor_id;type:varchar(128);not null;default:'';index;comment:锚定内容块ID" json:"anchorId"`
	Quote       string `gorm:"column:quote;type:varchar(1024);not null;default:'';comment:被选中文本快照" json:"quote"`
	Content     string `gorm:"column:content;type:text;not null;comment:评论正文" json:"content"`
	Mentions    string `gorm:"column:mentions;type:varchar(1024);not null;default:'';comment:被@用户ID，逗号分隔" json:"mentions"`
	CreatedBy   string `gorm:"column:created_by;type:varchar(64);not null;default:'';index;comment:评论人" json:"createdBy"`
	Resolved    bool   `gorm:"column:resolved;type:tinyint(1);not null;default:0;comment:是否已解决" json:"resolved"`
	ResolvedBy  string `gorm:"column:resolved_by;type:varchar(64);not null;default:'';comment:解决人" json:"resolvedBy"`
	ResolvedAt  string `gorm:"column:resolved_at;type:varchar(32);not null;default:'';comment:解决时间 RFC3339" json:"resolvedAt"`
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间" json:"deletedAt"`
}

func (DocComment) TableName() string { return "doc_comments" }

// DocShareLink 文档分享链接。表 doc_share_links。
// scope：anyone（任何人持链接）/ org（组织内）/ invited（仅受邀人，需登录态命中权限表）。
// 链接本身不绕过权限系统：ResolveShareLink 只做「把匿名访问者映射为某个角色」的一次性提升，
// 后续所有读写仍走 CheckPermission 的同一套判定。
type DocShareLink struct {
	models.Model
	LinkID        string `gorm:"column:link_id;type:varchar(64);uniqueIndex;not null;comment:链接业务ID" json:"linkId"`
	Token         string `gorm:"column:token;type:varchar(64);uniqueIndex;not null;comment:URL 中的随机令牌" json:"token"`
	ResourceType  string `gorm:"column:resource_type;type:varchar(16);not null;index:idx_share_res,priority:1;comment:space|doc" json:"resourceType"`
	ResourceID    string `gorm:"column:resource_id;type:varchar(64);not null;index:idx_share_res,priority:2;comment:资源ID" json:"resourceId"`
	Scope         string `gorm:"column:scope;type:varchar(16);not null;default:'anyone';comment:anyone|org|invited" json:"scope"`
	Role          string `gorm:"column:role;type:varchar(16);not null;default:'viewer';comment:链接授予的角色" json:"role"`
	PasswordHash  string `gorm:"column:password_hash;type:varchar(128);not null;default:'';comment:访问密码哈希，空为无密码" json:"passwordHash"`
	ExpiresAt     string `gorm:"column:expires_at;type:varchar(32);not null;default:'';comment:过期时间 RFC3339，空为永不过期" json:"expiresAt"`
	AllowDownload bool   `gorm:"column:allow_download;type:tinyint(1);not null;default:1;comment:允许下载" json:"allowDownload"`
	AllowCopy     bool   `gorm:"column:allow_copy;type:tinyint(1);not null;default:1;comment:允许复制" json:"allowCopy"`
	AllowComment  bool   `gorm:"column:allow_comment;type:tinyint(1);not null;default:0;comment:允许匿名评论" json:"allowComment"`
	Watermark     bool   `gorm:"column:watermark;type:tinyint(1);not null;default:0;comment:显示访问者水印" json:"watermark"`
	ViewCount     int64  `gorm:"column:view_count;type:bigint;not null;default:0;comment:访问次数" json:"viewCount"`
	CreatedBy     string `gorm:"column:created_by;type:varchar(64);not null;default:'';comment:创建人" json:"createdBy"`
	Status        int    `gorm:"column:status;type:tinyint;not null;default:1;index;comment:1有效 2已撤销" json:"status"`
}

func (DocShareLink) TableName() string { return "doc_share_links" }

// DocBlock 富媒体块（多维表格 / 画板）的结构化数据。表 doc_blocks。
// 正文的富文本在 Yjs 里协同，而 sheet/board 这类结构化块单独以 JSON 存储：
// 它们有自己的行列/节点模型，不适合塞进 ProseMirror 文档树。
type DocBlock struct {
	models.Model
	BlockID     string `gorm:"column:block_id;type:varchar(64);uniqueIndex;not null;comment:块业务ID" json:"blockId"`
	DocID       string `gorm:"column:doc_id;type:varchar(64);not null;index:idx_block_doc;comment:所属文档" json:"docId"`
	Kind        string `gorm:"column:kind;type:varchar(16);not null;comment:sheet|board" json:"kind"`
	Title       string `gorm:"column:title;type:varchar(255);not null;default:'';comment:块标题" json:"title"`
	PayloadJSON string `gorm:"column:payload_json;type:longtext;comment:结构化数据 JSON" json:"payloadJson"`
	Version     int    `gorm:"column:version;type:int;not null;default:1;comment:块版本" json:"version"`
	UpdatedBy   string `gorm:"column:updated_by;type:varchar(64);not null;default:'';comment:最后修改人" json:"updatedBy"`
}

func (DocBlock) TableName() string { return "doc_blocks" }
