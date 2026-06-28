package fileseed

import (
	"beaver/app/file/file_models"
	"fmt"

	"gorm.io/gorm"
)

// InitDefaultFiles 初始化默认文件数据
func InitDefaultFiles(db *gorm.DB) error {
	defaultFiles := []file_models.FileModel{
		{
			OriginalName: "defaultUserFileName",
			Size:         60317,
			Path:         "image/user.png",
			Md5:          "a9de5548bef8c10b92428fff61275c72",
			Type:         "image",
			FileKey:      "a9de5548bef8c10b92428fff61275c72.png",
			Source:       file_models.LocalSource,
		},
		{
			OriginalName: "defaultGroupFileName",
			Size:         90310,
			Path:         "image/group.png",
			Md5:          "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Type:         "image",
			FileKey:      "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Source:       file_models.LocalSource,
		},
	}

	for _, file := range defaultFiles {
		var count int64
		if err := db.Model(&file_models.FileModel{}).Where("file_key = ?", file.FileKey).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&file).Error; err != nil {
				return fmt.Errorf("创建默认文件失败: %w", err)
			}
			fmt.Printf("创建默认文件成功: %s\n", file.FileKey)
		}
	}

	return nil
}
