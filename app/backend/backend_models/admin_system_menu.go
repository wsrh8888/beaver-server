package backend_models

import "beaver/common/models"

/**
 * @description: 存储系统中的菜单信息，包括菜单与角色之间的映射关系。
 */
type AdminSystemMenu struct {
	models.Model
	ParentID *uint  `json:"ParentID" gorm:"comment:父菜单ID"` // 父菜单ID
	Path     string `json:"path" gorm:"comment:路由path"`    // 路由path
	Name     string `json:"name" gorm:"comment:路由name"`    // 路由name
	Hidden   bool   `json:"hidden" gorm:"comment:是否在列表隐藏"` // 是否在列表隐藏
	Sort     int    `json:"sort" gorm:"comment:排序标记"`      // 排序标记
	Title    string `json:"title" gorm:"comment:菜单名"`      // 菜单名
	Icon     string `json:"icon" gorm:"comment:菜单图标"`      // 菜单图标
}
