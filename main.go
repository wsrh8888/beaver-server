package main

import (
	"beaver/database"
	"flag"
	"fmt"
)

type Options struct {
	DB bool
}

// go run main.go -db
func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()

	if !opt.DB {
		return
	}

	dbMap, err := database.RunMigrations()
	if err != nil {
		fmt.Println(err)
		return
	}

<<<<<<< HEAD
	migrations := []struct {
		name string
		dsn  string
		run  func(*gorm.DB) error
	}{
		{
			name: "beaver_user",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_user?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(&user_models.UserModel{}, &user_models.UserChangeLogModel{})
			},
		},
		{
			name: "beaver_auth",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_auth?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(&auth_models.AuthCredentialModel{}, &auth_models.AuthDeviceModel{})
			},
		},
		{
			name: "beaver_friend",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_friend?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&friend_models.FriendModel{},
					&friend_models.FriendVerifyModel{},
					&friend_models.FriendBlockModel{},
				)
			},
		},
		{
			name: "beaver_group",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_group?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&group_models.GroupModel{},
					&group_models.GroupMemberModel{},
					&group_models.GroupJoinRequestModel{},
					&group_models.GroupMemberChangeLogModel{},
					&group_models.GroupBotModel{},
				)
			},
		},
		{
			name: "beaver_chat",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_chat?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&chat_models.ChatMessage{},
					&chat_models.ChatConversationMeta{},
					&chat_models.ChatUserConversation{},
					&chat_models.ChatUserDelete{},
					&chat_models.ChatForward{},
				)
			},
		},
		{
			name: "beaver_moment",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_moment?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&moment_models.MomentModel{},
					&moment_models.MomentLikeModel{},
					&moment_models.MomentCommentModel{},
				)
			},
		},
		{
			name: "beaver_emoji",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_emoji?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&emoji_models.EmojiPackage{},
					&emoji_models.Emoji{},
					&emoji_models.EmojiPackageEmoji{},
					&emoji_models.EmojiPackageCollect{},
					&emoji_models.EmojiCollectEmoji{},
				)
			},
		},
		{
			name: "beaver_file",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_file?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(&file_models.FileModel{})
			},
		},
		{
			name: "beaver_notification",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_notification?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&notification_models.NotificationEvent{},
					&notification_models.NotificationInbox{},
					&notification_models.NotificationRead{},
					&notification_models.PushRegistrationModel{},
				)
			},
		},
		{
			name: "beaver_call",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_call?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(&call_models.CallSession{}, &call_models.CallParticipant{})
			},
		},
		{
			name: "beaver_open",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_open?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
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
				)
			},
		},
		{
			name: "beaver_platform",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_platform?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&platform_models.TrackBucket{},
					&platform_models.TrackEvent{},
					&platform_models.TrackLogger{},
					&platform_models.FeedbackModel{},
					&platform_models.ContentReportModel{},
					&platform_models.UpdateApp{},
					&platform_models.UpdateArchitecture{},
					&platform_models.UpdateVersion{},
					&platform_models.UpdateStrategy{},
					&platform_models.UpdateReport{},
				)
			},
		},
		{
			name: "beaver_backend",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_backend?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&backend_models.AdminUser{},
					&backend_models.AdminSystemAuthority{},
					&backend_models.AdminSystemAuthorityMenu{},
					&backend_models.AdminSystemAuthorityUser{},
					&backend_models.AdminSystemMenu{},
					&backend_models.AdminModerationCase{},
					&backend_models.AdminOperationLog{},
					&backend_models.AdminSensitiveWord{},
				)
			},
		},
		{
			name: "beaver_datasync",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_datasync?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(&datasync_models.DatasyncModel{})
			},
		},
	}

	for _, m := range migrations {
		db := coregorm.InitGorm(m.dsn)
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		if err := m.run(db); err != nil {
			fmt.Printf("%s 表结构生成失败: %v\n", m.name, err)
			return
		}
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		fmt.Printf("%s 表结构生成成功\n", m.name)
	}

	fileDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_file?charset=utf8mb4&parseTime=True&loc=Local")
	platformDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_platform?charset=utf8mb4&parseTime=True&loc=Local")
	userDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_user?charset=utf8mb4&parseTime=True&loc=Local")
	authDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_auth?charset=utf8mb4&parseTime=True&loc=Local")
	openDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_open?charset=utf8mb4&parseTime=True&loc=Local")
	_ = database.InitFileData(fileDB)
	_ = database.InitUpdateApp(platformDB)
	_ = database.InitUpdateStrategy(platformDB)
	if err := database.InitDefaultUser(userDB, authDB, openDB); err != nil {
=======
	if err := database.InitDefaultUser(
		dbMap["beaver_user"],
		dbMap["beaver_auth"],
		dbMap["beaver_open"],
	); err != nil {
>>>>>>> 5b160f761511999451c7965cce0863719701c988
		fmt.Printf("默认用户初始化失败: %v\n", err)
		return
	}

	fmt.Println("所有库表结构生成完成")
}
