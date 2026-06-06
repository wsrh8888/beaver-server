package database

import (
	"beaver/app/auth/auth_models"
	"beaver/app/backend/backend_models"
	"beaver/app/call/call_models"
	"beaver/app/chat/chat_models"
	"beaver/app/datasync/datasync_models"
	"beaver/app/emoji/emoji_models"
	"beaver/app/file/file_models"
	"beaver/app/friend/friend_models"
	"beaver/app/group/group_models"
	"beaver/app/moment/moment_models"
	"beaver/app/notification/notification_models"
	"beaver/app/open/open_models"
	"beaver/app/platform/platform_models"
	"beaver/app/user/user_models"
	"beaver/core/coregorm"
	"fmt"

	"gorm.io/gorm"
)

const (
	mysqlUser = "root:123456"
	mysqlAddr = "127.0.0.1:3306"
)

// Migration 单库迁移配置：库名 + 表模型 + 可选种子数据
type Migration struct {
	Name   string
	Models []any
	Seed   func(*gorm.DB) error
}

// DSN 根据库名生成连接串，dbName 为空时连接实例（用于 CREATE DATABASE）
func DSN(dbName string) string {
	if dbName == "" {
		return fmt.Sprintf("%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlAddr)
	}
	return fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlAddr, dbName)
}

// AllMigrations 全部业务库迁移定义（新增库只改这一处）
func AllMigrations() []Migration {
	return []Migration{
		{
			Name: "beaver_user",
			Models: []any{
				&user_models.UserModel{},
				&user_models.UserChangeLogModel{},
			},
		},
		{
			Name: "beaver_auth",
			Models: []any{
				&auth_models.AuthCredentialModel{},
				&auth_models.AuthDeviceModel{},
			},
		},
		{
			Name: "beaver_friend",
			Models: []any{
				&friend_models.FriendModel{},
				&friend_models.FriendVerifyModel{},
				&friend_models.FriendBlockModel{},
			},
		},
		{
			Name: "beaver_group",
			Models: []any{
				&group_models.GroupModel{},
				&group_models.GroupMemberModel{},
				&group_models.GroupJoinRequestModel{},
				&group_models.GroupMemberChangeLogModel{},
				&group_models.GroupBotModel{},
			},
		},
		{
			Name: "beaver_chat",
			Models: []any{
				&chat_models.ChatMessage{},
				&chat_models.ChatConversationMeta{},
				&chat_models.ChatUserConversation{},
				&chat_models.ChatUserDelete{},
				&chat_models.ChatForward{},
			},
		},
		{
			Name: "beaver_moment",
			Models: []any{
				&moment_models.MomentModel{},
				&moment_models.MomentLikeModel{},
				&moment_models.MomentCommentModel{},
			},
		},
		{
			Name: "beaver_emoji",
			Models: []any{
				&emoji_models.EmojiPackage{},
				&emoji_models.Emoji{},
				&emoji_models.EmojiPackageEmoji{},
				&emoji_models.EmojiPackageCollect{},
				&emoji_models.EmojiCollectEmoji{},
			},
		},
		{
			Name:   "beaver_file",
			Models: []any{&file_models.FileModel{}},
			Seed:   InitFileData,
		},
		{
			Name: "beaver_notification",
			Models: []any{
				&notification_models.NotificationEvent{},
				&notification_models.NotificationInbox{},
				&notification_models.NotificationRead{},
				&notification_models.PushRegistrationModel{},
			},
		},
		{
			Name: "beaver_call",
			Models: []any{
				&call_models.CallSession{},
				&call_models.CallParticipant{},
			},
		},
		{
			Name: "beaver_open",
			Models: []any{
				&open_models.OpenDeveloper{},
				&open_models.OpenApp{},
				&open_models.OpenAppOAuth{},
				&open_models.OpenAppRobot{},
				&open_models.OpenAppSecurity{},
				&open_models.OpenAppEventSubscription{},
				&open_models.OpenOAuthToken{},
				&open_models.OpenOAuthCode{},
				&open_models.OpenOAuthQrCode{},
				&open_models.OpenBotModel{},
				&open_models.OpenWebhookLog{},
				&open_models.OpenRobotSendLog{},
			},
		},
		{
			Name: "beaver_platform",
			Models: []any{
				&platform_models.TrackBucket{},
				&platform_models.TrackEvent{},
				&platform_models.TrackLogger{},
				&platform_models.FeedbackModel{},
				&platform_models.UpdateApp{},
				&platform_models.UpdateArchitecture{},
				&platform_models.UpdateVersion{},
				&platform_models.UpdateStrategy{},
				&platform_models.UpdateReport{},
			},
			Seed: seedPlatform,
		},
		{
			Name: "beaver_backend",
			Models: []any{
				&backend_models.AdminUser{},
				&backend_models.AdminSystemAuthority{},
				&backend_models.AdminSystemAuthorityMenu{},
				&backend_models.AdminSystemAuthorityUser{},
				&backend_models.AdminSystemMenu{},
			},
		},
		{
			Name:   "beaver_datasync",
			Models: []any{&datasync_models.DatasyncModel{}},
		},
	}
}

func seedPlatform(db *gorm.DB) error {
	if err := InitUpdateApp(db); err != nil {
		return err
	}
	return InitUpdateStrategy(db)
}

// RunMigrations 建库、迁表、单库种子数据，返回各库连接供跨库初始化使用
func RunMigrations() (map[string]*gorm.DB, error) {
	migrations := AllMigrations()
	serverDB := coregorm.InitGorm(DSN(""))

	for _, m := range migrations {
		sql := fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
			m.Name,
		)
		if err := serverDB.Exec(sql).Error; err != nil {
			return nil, fmt.Errorf("创建数据库 %s 失败: %w", m.Name, err)
		}
	}

	dbMap := make(map[string]*gorm.DB, len(migrations))
	for _, m := range migrations {
		db := coregorm.InitGorm(DSN(m.Name))
		dbMap[m.Name] = db

		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		if err := db.AutoMigrate(m.Models...); err != nil {
			return nil, fmt.Errorf("%s 表结构生成失败: %w", m.Name, err)
		}
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		fmt.Printf("%s 表结构生成成功\n", m.Name)

		if m.Seed != nil {
			if err := m.Seed(db); err != nil {
				return nil, fmt.Errorf("%s 种子数据初始化失败: %w", m.Name, err)
			}
		}
	}

	return dbMap, nil
}
