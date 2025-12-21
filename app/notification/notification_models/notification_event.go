package notification_models

import (
	"beaver/common/models"

	"gorm.io/datatypes"
)

// NotificationEvent 事件主表：用于幂等和审计，跨服务统一存储事件元数据
type NotificationEvent struct {
	models.Model
	EventID    string         `gorm:"size:64;uniqueIndex;not null" json:"eventId"` // 全局事件ID（雪花/ULID）
	EventType  string         `gorm:"size:32;index;not null" json:"eventType"`     // 事件类型：friend_request/moment_like/...
	Category   string         `gorm:"size:32;index;not null" json:"category"`      // 场景分类：social/system/group/moment等
	Version    int64          `gorm:"not null;default:0;index" json:"version"`     // 全局递增版本号（建议取雪花时间或全局序列），用于增量同步/纠偏
	FromUserID *string        `gorm:"size:64;index" json:"fromUserId,omitempty"`   // 事件触发方
	TargetID   *string        `gorm:"size:64;index" json:"targetId,omitempty"`     // 目标对象（如动态ID、群ID）
	TargetType string         `gorm:"size:32;index" json:"targetType"`             // 目标类型：moment/group/user/message等
	Payload    datatypes.JSON `json:"payload"`                                     // 事件扩展数据，前端可直接渲染
	Priority   int8           `gorm:"not null;default:5;index" json:"priority"`    // 优先级（1最高，9最低），便于队列/推送调度
	Status     int8           `gorm:"not null;default:1;index" json:"status"`      // 1=有效 2=撤回/隐藏 3=失效
	DedupHash  string         `gorm:"size:128;index" json:"dedupHash"`             // 去重哈希，用于点赞合并等
}
