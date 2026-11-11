package doc_models

import "beaver/common/models"

// DocAuditLog 文档域操作审计。表 doc_audit_logs。
// 所有写操作在 doc_rpc 内统一落一条，detail 为该操作的上下文 JSON（不含正文，避免审计表膨胀）。
type DocAuditLog struct {
	models.Model
	LogID        string `gorm:"column:log_id;type:varchar(64);uniqueIndex;not null;comment:日志业务ID" json:"logId"`
	ActorUserID  string `gorm:"column:actor_user_id;type:varchar(64);not null;default:'';index;comment:操作人" json:"actorUserId"`
	Action       string `gorm:"column:action;type:varchar(64);not null;index;comment:操作类型" json:"action"`
	ResourceType string `gorm:"column:resource_type;type:varchar(16);not null;default:'';comment:资源类型" json:"resourceType"`
	ResourceID   string `gorm:"column:resource_id;type:varchar(64);not null;default:'';index;comment:资源ID" json:"resourceId"`
	Detail       string `gorm:"column:detail;type:varchar(2048);not null;default:'';comment:上下文 JSON" json:"detail"`
	CreatedAtStr string `gorm:"column:created_at_str;type:varchar(32);not null;default:'';index;comment:发生时间 RFC3339" json:"createdAtStr"`
}

func (DocAuditLog) TableName() string { return "doc_audit_logs" }

// DocNotification 文档域通知中心。表 doc_notifications。
// 覆盖 @提及、评论回复、被分享、权限变更、版本回滚等事件；客户端在文档工作台内消费。
type DocNotification struct {
	models.Model
	NotificationID string `gorm:"column:notification_id;type:varchar(64);uniqueIndex;not null;comment:通知业务ID" json:"notificationId"`
	UserID         string `gorm:"column:user_id;type:varchar(64);not null;index:idx_notify_user,priority:1;comment:接收人" json:"userId"`
	DocID          string `gorm:"column:doc_id;type:varchar(64);not null;default:'';index;comment:相关文档" json:"docId"`
	Kind           string `gorm:"column:kind;type:varchar(32);not null;comment:mention|comment_reply|share|permission|rollback" json:"kind"`
	ActorUserID    string `gorm:"column:actor_user_id;type:varchar(64);not null;default:'';comment:触发人" json:"actorUserId"`
	Title          string `gorm:"column:title;type:varchar(255);not null;default:'';comment:通知标题" json:"title"`
	Excerpt        string `gorm:"column:excerpt;type:varchar(512);not null;default:'';comment:摘要" json:"excerpt"`
	IsRead         bool   `gorm:"column:is_read;type:tinyint(1);not null;default:0;index:idx_notify_user,priority:2;comment:是否已读" json:"isRead"`
	CreatedAtStr   string `gorm:"column:created_at_str;type:varchar(32);not null;default:'';index;comment:发生时间 RFC3339" json:"createdAtStr"`
}

func (DocNotification) TableName() string { return "doc_notifications" }

// DocIndexOutbox 知识库索引事务外盒。表 doc_index_outbox。
// 正文每次落库都在同一事务写一条 pending，平台知识库按 (doc_id, version) 拉取 GetDocContent 建索引；
// 这样 doc 域不需要反向依赖 platform_rpc，也不会出现「正文已改但索引没跟上」的丢事件。
type DocIndexOutbox struct {
	models.Model
	DocID     string `gorm:"column:doc_id;type:varchar(64);not null;uniqueIndex:idx_outbox_doc_ver,priority:1;comment:文档ID" json:"docId"`
	Version   int    `gorm:"column:version;type:int;not null;uniqueIndex:idx_outbox_doc_ver,priority:2;comment:正文版本" json:"version"`
	Digest    string `gorm:"column:digest;type:varchar(64);not null;default:'';comment:正文哈希" json:"digest"`
	Status    string `gorm:"column:status;type:varchar(16);not null;default:'pending';index;comment:pending|done|failed" json:"status"`
	Attempts  int    `gorm:"column:attempts;type:int;not null;default:0;comment:处理次数" json:"attempts"`
	CreatedStr string `gorm:"column:created_str;type:varchar(32);not null;default:'';index;comment:入队时间 RFC3339" json:"createdStr"`
}

func (DocIndexOutbox) TableName() string { return "doc_index_outbox" }
