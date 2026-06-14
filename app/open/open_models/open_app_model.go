package open_models

import (
	"time"

	"gorm.io/gorm"
)

// ==================== OpenApp 应用主表 ====================

// OpenApp 开放平台应用主表（对标钉钉开放平台）
type OpenApp struct {
	gorm.Model
	// 身份认证
	AppID     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用唯一标识"`
	AppSecret string `gorm:"type:varchar(128);not null;comment:应用密钥"`
	// 基础信息
	Name        string `gorm:"type:varchar(100);not null;comment:应用名称"`
	Description string `gorm:"type:text;comment:应用描述"`
	Icon        string `gorm:"type:varchar(500);comment:应用图标URL"`
	OwnerUserID string `gorm:"type:varchar(64);index;comment:所属用户ID"`
	// 应用类型与分类
	AppType  int    `gorm:"type:tinyint;default:1;comment:应用类型 1-自建应用 2-第三方应用"`
	Category string `gorm:"type:varchar(50);comment:应用分类 office/entertainment/tool"`
	// 状态管理
	Status      int        `gorm:"type:tinyint;default:0;comment:状态 0草稿 1已发布 2禁用"`
	AuditStatus int        `gorm:"type:tinyint;default:0;comment:审核状态 0待审核 1已通过 2已拒绝"`
	AuditedBy   string     `gorm:"type:varchar(64);comment:审核人ID"`
	AuditedAt   *time.Time `gorm:"type:datetime;comment:审核时间"`
	// 能力开关（快速筛选）
	EnableRobot   int `gorm:"type:tinyint;default:0;comment:是否启用智能机器人能力 1是 0否"`
	EnableOAuth   int `gorm:"type:tinyint;default:0;comment:是否启用OAuth能力 1是 0否"`
	EnableWebhook int `gorm:"type:tinyint;default:0;comment:是否启用Webhook能力 1是 0否"`
	// 客户端标识（用于 JSSDK 鉴权和自定义协议）
	AgentId string `gorm:"type:varchar(64);index;comment:微应用ID(用于JSSDK鉴权)"`
	Scheme  string `gorm:"type:varchar(64);comment:客户端回调协议(Scheme)，如 beaver://"`
	// 审计字段
	LastModifiedBy string `gorm:"type:varchar(64);comment:最后修改人ID"`
	Version        int    `gorm:"default:1;comment:配置版本号"`
}
