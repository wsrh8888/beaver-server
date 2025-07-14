package main

import (
	"beaver/app/chat/chat_models"
	"beaver/app/emoji/emoji_models"
	"beaver/app/feedback/feedback_models"
	"beaver/app/file/file_models"
	"beaver/app/friend/friend_models"
	"beaver/app/group/group_models"
	"beaver/app/user/user_models"
	"beaver/core"
	"flag"
	"fmt"
)

type Options struct {
	DB bool
}

// go run main.go  -db
func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()

	if opt.DB {
		db := core.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver?charset=utf8mb4&parseTime=True&loc=Local")

		// 禁用外键检查
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

		// 创建所有表
		err := db.AutoMigrate(
			// 基础表
			&user_models.UserModel{},
			&friend_models.FriendModel{},
			&friend_models.FriendVerifyModel{},
			&chat_models.ChatModel{},
			&chat_models.ChatUserConversationModel{},
			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&file_models.FileModel{},
			&emoji_models.EmojiPackage{},
			&emoji_models.Emoji{},
			&emoji_models.EmojiPackageEmoji{},
			&emoji_models.EmojiPackageCollect{},
			&emoji_models.EmojiCollectEmoji{},
			&feedback_models.FeedbackModel{},
		)

		if err != nil {
			fmt.Printf("表结构生成失败: %v\n", err)
			return
		}

		fmt.Println("所有表结构生成成功")
	}
}
