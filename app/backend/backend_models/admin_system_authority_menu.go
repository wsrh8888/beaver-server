package backend_models

import "beaver/common/models"

/**
 * @description: 表示角色-菜单关联关系数据模型

 */
type AdminSystemAuthorityMenu struct {
	models.Model
	AuthorityID uint `json:"authorityId" gorm:"comment:角色ID;index:idx_authority_menu,unique"`
	MenuID      uint `json:"menuId" gorm:"comment:菜单ID;index:idx_authority_menu,unique"`
}
