package backend_models

import "beaver/common/models"

/**
 * @description: 表示用户-角色关联关系数据模型
 * 注意：微服务架构中不使用数据库外键约束，数据一致性由应用层保证
 */
type AdminSystemAuthorityUser struct {
	models.Model
	UserID      string `json:"userId" gorm:"size:64;comment:用户ID;index:idx_user_authority,unique"`
	AuthorityID uint   `json:"authorityId" gorm:"comment:角色ID;index:idx_user_authority,unique"`
}
