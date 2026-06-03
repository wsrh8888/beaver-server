package main

import (
	"beaver/app/auth/auth_models"
	"beaver/app/backend/backend_models"
	"beaver/app/call/call_models"
	"beaver/app/chat/chat_models"
	"beaver/app/datasync/datasync_models"
	"beaver/app/emoji/emoji_models"
	"beaver/app/feedback/feedback_models"
	"beaver/app/file/file_models"
	"beaver/app/friend/friend_models"
	"beaver/app/group/group_models"
	"beaver/app/moment/moment_models"
	"beaver/app/notification/notification_models"
	"beaver/app/open/open_models"
	"beaver/app/track/track_models"
	"beaver/app/platform/platform_models"
	"beaver/app/user/user_models"
	"beaver/core/coregorm"
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
		db := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver?charset=utf8mb4&parseTime=True&loc=Local")

		// 禁用外键检查
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

		// 创建所有表
		err := db.AutoMigrate(
			// 基础表
			&user_models.UserModel{},
			&user_models.UserChangeLogModel{},

			// 认证相关表
			&auth_models.AuthCredentialModel{},
			&auth_models.AuthDeviceModel{},

			// 好友相关表
			&friend_models.FriendModel{},
			&friend_models.FriendVerifyModel{},

			// 聊天相关表
			&chat_models.ChatMessage{},
			&chat_models.ChatConversationMeta{},
			&chat_models.ChatUserConversation{},
			&chat_models.ChatUserDelete{},
			&chat_models.ChatForward{},

			// 文件相关表
			&file_models.FileModel{},

			// 动态相关表
			&moment_models.MomentModel{},
			&moment_models.MomentLikeModel{},
			&moment_models.MomentCommentModel{},

			// 表情相关表
			&emoji_models.EmojiPackage{},
			&emoji_models.Emoji{},
			&emoji_models.EmojiPackageEmoji{},
			&emoji_models.EmojiPackageCollect{},
			&emoji_models.EmojiCollectEmoji{},

			// 埋点相关的
			&track_models.TrackBucket{},
			&track_models.TrackEvent{},
			&track_models.TrackLogger{},

			// 反馈相关表
			&feedback_models.FeedbackModel{},

			// 版本管理相关表
			&platform_models.UpdateApp{},
			&platform_models.UpdateArchitecture{},
			&platform_models.UpdateVersion{},
			&platform_models.UpdateStrategy{},
			&platform_models.UpdateReport{},

			&datasync_models.DatasyncModel{},

			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&group_models.GroupJoinRequestModel{},
			&group_models.GroupMemberChangeLogModel{},
			&group_models.GroupBotModel{},

			// 通知中心表
			&notification_models.NotificationEvent{},
			&notification_models.NotificationInbox{},
			&notification_models.NotificationRead{},
			&notification_models.PushRegistrationModel{},

			// 音视频通话表
			&call_models.CallSession{},
			&call_models.CallParticipant{},

			// 后台管理相关表
			&backend_models.AdminUser{},
			&backend_models.AdminSystemAuthority{},
			&backend_models.AdminSystemAuthorityMenu{},
			&backend_models.AdminSystemAuthorityUser{},
			&backend_models.AdminSystemMenu{},

			// 开放平台相关表
			&open_models.OpenDeveloper{},
			&open_models.OpenApp{},
			&open_models.OpenAppOAuth{},             // OAuth 配置表
			&open_models.OpenAppRobot{},             // Robot 配置表
			&open_models.OpenAppSecurity{},          // Security 配置表
			&open_models.OpenAppEventSubscription{}, // 事件订阅表
			&open_models.OpenOAuthToken{},
			&open_models.OpenOAuthCode{},
			&open_models.OpenOAuthQrCode{},
			&open_models.OpenBotModel{},
		)

		if err != nil {
			fmt.Printf("表结构生成失败: %v\n", err)
			return
		}

		fmt.Println("所有表结构生成成功")
		database.InitAllData(db)
	}
}
