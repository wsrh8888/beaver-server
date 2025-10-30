package backend_models

import (
	"beaver/common/models"
)

/**
 * @description:  存储系统中的角色信息以及角色与用户、菜单之间的映射关系。
 */
type AdminSystemAuthority struct {
	models.Model
	Name        string `json:"authorityName"` // 角色名
	Description string `json:"description"`   // 角色描述
	// MenuModel   []AdminAuthorityMenu `gorm:"foreignkey:Id;references:AuthorityID" json:"-"`
	// AdminUserModel []AdminAuthorityUser `gorm:"foreignkey:Id;references:AuthorityID" json:"-"`
}
