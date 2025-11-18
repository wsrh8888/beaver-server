package backend_models

import (
	"beaver/common/models"
)

/**
 * @description: 存储系统中的角色信息
 * 注意：微服务架构中不使用数据库外键约束，关联关系通过关联表维护
 */
type AdminSystemAuthority struct {
	models.Model
	Name        string `json:"authorityName" gorm:"size:32;unique;index;comment:角色名"` // 角色名（唯一）
	Description string `json:"description" gorm:"size:256;comment:角色描述"`              // 角色描述
	Status      int8   `json:"status" gorm:"default:1;index;comment:状态"`              // 1:启用 2:禁用
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`                      // 排序
}
