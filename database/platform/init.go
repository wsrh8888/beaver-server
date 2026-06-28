package platform

import (
	"beaver/database/platform/track"
	"beaver/database/platform/update"

	"gorm.io/gorm"
)

// InitPlatform 初始化 beaver_platform 库默认数据
func InitPlatform(db *gorm.DB) error {
	if err := update.InitUpdateApp(db); err != nil {
		return err
	}
	return track.InitBuckets(db)
}
