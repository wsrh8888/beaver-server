package database

import (
	"beaver/app/file/file_models"
	"log"

	"gorm.io/gorm"
)

func InitFileData(db *gorm.DB) error {
	// 初始化默认文件数据
	defaultFiles := []file_models.FileModel{
		{
			OriginalName: "defaultUserFileName",
			Size:         60317,
			Path:         "image/user.png",
			Md5:          "a9de5548bef8c10b92428fff61275c72",
			Type:         "image",
			FileName:     "a9de5548bef8c10b92428fff61275c72.png",
			Source:       file_models.LocalSource,
		},
		{
			OriginalName: "defaultGroupFileName",
			Size:         90310,
			Path:         "image/group.png",
			Md5:          "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Type:         "image",
			FileName:     "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Source:       file_models.LocalSource,
		},
	}

	for _, file := range defaultFiles {
		var count int64
		if err := db.Model(&file_models.FileModel{}).Where("file_name = ?", file.FileName).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&file).Error; err != nil {
				log.Printf("创建默认文件失败: %v", err)
				return err
			} else {
				log.Printf("创建默认文件成功: %s", file.FileName)
			}
		}
	}

	return nil
}
