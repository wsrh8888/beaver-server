package database

import (
	"beaver/app/file/file_models"
	"fmt"

	"gorm.io/gorm"
)

// 初始化文件表数据
func InitFileData(db *gorm.DB) error {
	// 定义默认文件数据
	defaultFiles := []file_models.FileModel{
		{
			FileName: "defaultUserAvatar",
			Size:     60317,
			Path:     "image/872f5fe419791dbf31c6635e8a7a6594.png",
			Hash:     "872f5fe419791dbf31c6635e8a7a6594",
			Type:     "image",
			FileID:   "eb3dad2d-4b7f-44c2-9af5-50ad9f76ff81",
		},
		{
			FileName: "defaultGroupAvatar",
			Size:     90310,
			Path:     "image/83f8cc2c12a508444281b3181c62608f.png",
			Hash:     "83f8cc2c12a508444281b3181c62608f",
			Type:     "image",
			FileID:   "71e4be6c-b477-4fce-8348-9cc53349ef28",
		},
	}

	// 对每条记录检查是否存在，不存在才插入
	for _, file := range defaultFiles {
		var count int64
		// 检查记录是否存在
		db.Model(&file_models.FileModel{}).Where("file_id = ?", file.FileID).Count(&count)

		if count == 0 {
			// 记录不存在，创建它
			if err := db.Create(&file).Error; err != nil {
				return fmt.Errorf("创建文件记录失败 (ID: %s): %v", file.FileID, err)
			}
			fmt.Printf("已创建文件记录: %s\n", file.FileID)
		} else {
			fmt.Printf("文件记录已存在，跳过: %s\n", file.FileID)
		}
	}
	return nil
}
