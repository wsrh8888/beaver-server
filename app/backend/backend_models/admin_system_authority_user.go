package backend_models

import "beaver/common/models"

/**
 * @description: 表示用户-角色关联关系数据模型
 */
type AdminSystemAuthorityUser struct {
	models.Model
	UserID         string               `json:"userId" gorm:"comment:用户ID"`
	AuthorityID    uint                 `json:"authorityId" gorm:"comment:角色ID"`
	AuthorityModel AdminSystemAuthority `gorm:"foreignKey:AuthorityID;references:Id" json:"-"`
}
