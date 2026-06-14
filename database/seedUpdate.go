package database

import "gorm.io/gorm"

// SeedUpdateData 幂等初始化升级默认应用
func SeedUpdateData(db *gorm.DB) error {
	return InitUpdateApp(db)
}
