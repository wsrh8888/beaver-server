package backend_models

import "beaver/common/models"

/**
 * @description: 存储系统中的菜单信息
 * 注意：菜单与角色的关联关系通过 AdminSystemAuthorityMenu 表维护
 */
type AdminSystemMenu struct {
	models.Model
	ParentID   *uint  `json:"parentId" gorm:"index;comment:父菜单ID"`             // 父菜单ID（NULL表示顶级菜单）
	Path       string `json:"path" gorm:"size:128;index;comment:路由path"`       // 路由path
	Name       string `json:"name" gorm:"size:64;unique;index;comment:路由name"` // 路由name（唯一，用于路由匹配）
	Component  string `json:"component" gorm:"size:256;comment:组件路径"`          // 前端组件路径（如：views/user/index.vue）
	Hidden     bool   `json:"hidden" gorm:"default:false;comment:是否在列表隐藏"`     // 是否在列表隐藏
	Sort       int    `json:"sort" gorm:"default:0;index;comment:排序标记"`        // 排序标记
	Title      string `json:"title" gorm:"size:32;comment:菜单名"`                // 菜单名
	Icon       string `json:"icon" gorm:"size:64;comment:菜单图标"`                // 菜单图标
	Permission string `json:"permission" gorm:"size:64;index;comment:权限标识"`    // 权限标识（用于按钮级权限控制，如：user:create）
	Status     int8   `json:"status" gorm:"default:1;index;comment:状态"`        // 1:启用 2:禁用
}
