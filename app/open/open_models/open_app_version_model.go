package open_models

import (
	"gorm.io/gorm"
)

// OpenAppVersion 应用版本历史表
type OpenAppVersion struct {
	gorm.Model
	AppID        string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Version      string `gorm:"type:varchar(20);not null;comment:版本号"`
	Description  string `gorm:"type:text;comment:版本说明"`
	Visibility   string `gorm:"type:varchar(20);default:self;comment:可见范围 self/partial/all"`
	Status       string `gorm:"type:varchar(20);default:draft;comment:状态 draft/reviewing/approved/rejected/published"`
	Capabilities string `gorm:"type:text;comment:该版本包含的能力(JSON数组)"`
	CreatedBy    string `gorm:"type:varchar(64);comment:创建人"`
}
