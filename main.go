package main

import (
	"beaver/app/backend/backend_models"
	"beaver/app/chat/chat_models"
	"beaver/app/datasync/datasync_models"
	"beaver/app/emoji/emoji_models"
	"beaver/app/feedback/feedback_models"
	"beaver/app/file/file_models"
	"beaver/app/friend/friend_models"
	"beaver/app/group/group_models"
	"beaver/app/moment/moment_models"
	"beaver/app/notification/notification_models"
	"beaver/app/track/track_models"
	"beaver/app/update/update_models"
	"beaver/app/user/user_models"
	"beaver/core"
	"beaver/database"
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
			&user_models.UserChangeLogModel{},
			&friend_models.FriendModel{},
			&friend_models.FriendVerifyModel{},
			&chat_models.ChatMessage{},
			&chat_models.ChatConversationMeta{},
			&chat_models.ChatUserConversation{},

			&file_models.FileModel{},
			&moment_models.MomentModel{},
			&moment_models.MomentLikeModel{},
			&moment_models.MomentCommentModel{},
			&emoji_models.EmojiPackage{},
			&emoji_models.Emoji{},
			&emoji_models.EmojiPackageEmoji{},
			&emoji_models.EmojiPackageCollect{},
			&emoji_models.EmojiCollectEmoji{},

			&track_models.TrackBucket{},
			&track_models.TrackEvent{},
			&track_models.TrackLogger{},
			&feedback_models.FeedbackModel{},

			// 版本管理相关表
			&update_models.UpdateApp{},
			&update_models.UpdateArchitecture{},
			&update_models.UpdateVersion{},
			&update_models.UpdateStrategy{},
			&update_models.UpdateReport{},

			&datasync_models.DatasyncModel{},

			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&group_models.GroupJoinRequestModel{},
			&group_models.GroupMemberChangeLogModel{},

			// 通知中心表
			&notification_models.NotificationEvent{},
			&notification_models.NotificationInbox{},
			&notification_models.NotificationReadCursor{},

			// 后台管理相关表
			&backend_models.AdminUser{},
			&backend_models.AdminSystemAuthority{},
			&backend_models.AdminSystemAuthorityMenu{},
			&backend_models.AdminSystemAuthorityUser{},
			&backend_models.AdminSystemMenu{},
		)

		if err != nil {
			fmt.Printf("表结构生成失败: %v\n", err)
			return
		}

		fmt.Println("所有表结构生成成功")
		database.InitAllData(db)
	}
}
