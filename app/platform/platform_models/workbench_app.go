package platform_models

import "beaver/common/models"

// WorkbenchApp IM 运营工作台应用（后台配置、全员展示，客户端 WebView 打开）
type WorkbenchApp struct {
	models.Model
	WorkbenchAppID string `gorm:"size:64;uniqueIndex;not null;comment:应用业务ID" json:"workbenchAppId"`
	Name           string `gorm:"size:64;not null;comment:应用名称" json:"name"`
	Description    string `gorm:"size:500;comment:应用描述" json:"description"`
	Icon           string `gorm:"size:500;comment:应用图标URL" json:"icon"`
	EntryURL       string `gorm:"size:1000;not null;comment:入口URL(WebView加载地址)" json:"entryUrl"`
	Category       string `gorm:"size:32;index;comment:分组/分类(展示用)" json:"category"`
	Sort           int    `gorm:"not null;default:0;index;comment:排序(越小越靠前)" json:"sort"`
	Status         int8   `gorm:"type:tinyint;not null;default:0;index;comment:状态 0下架 1上架" json:"status"`
	CreatedBy      string             `gorm:"size:64;comment:创建人(管理员ID)" json:"createdBy"`
	LastModifiedBy string             `gorm:"size:64;comment:最后修改人(管理员ID)" json:"lastModifiedBy"`
	Remark         string             `gorm:"size:500;comment:运营备注(不对客户端暴露)" json:"remark"`
}

func (WorkbenchApp) TableName() string {
	return "workbench_apps"
}
