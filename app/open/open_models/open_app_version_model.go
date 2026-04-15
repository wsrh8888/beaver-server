package open_models

import (
	"gorm.io/gorm"
)

// OpenAppVersion 应用版本历史表
type OpenAppVersion struct {
	gorm.Model
	AppID     string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Version   string `gorm:"type:varchar(20);not null;comment:版本号"`
	ChangeLog string `gorm:"type:text;comment:更新日志"`
	Status    int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
	CreatedBy string `gorm:"type:varchar(64);comment:创建人"`
}
