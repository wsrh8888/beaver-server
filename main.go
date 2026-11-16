/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package main

import (
	auth_models "beaver/app/account/account_models/auth"
	"beaver/app/message/message_models"
	"beaver/app/datasync/datasync_models"
	emoji_models "beaver/app/asset/asset_models/emoji"
	file_models "beaver/app/asset/asset_models/file"
	"beaver/app/relation/relation_models/friend"
	group_models "beaver/app/relation/relation_models/group"
	"beaver/app/notification/notification_models"
	"beaver/app/platform/platform_models/feedback"
	"beaver/app/platform/platform_models/report"
	"beaver/app/platform/platform_models/track"
	"beaver/app/platform/platform_models/update"
	"beaver/app/platform/platform_models/workbench"
	org_models "beaver/app/platform/platform_models/org"
	user_models "beaver/app/account/account_models/user"
	"beaver/core/coregorm"
	fileseed "beaver/database/file"
	"beaver/database/platform"
	userseed "beaver/database/user"
	"flag"
	"fmt"

	"gorm.io/gorm"
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

	serverDSN := "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	serverDB := coregorm.InitGorm(serverDSN)

	databases := []string{
		"beaver_account", "beaver_relation", "beaver_message",
		"beaver_asset", "beaver_notification",
		"beaver_platform", "beaver_datasync", "beaver_agent",
	}
	for _, name := range databases {
		if err := serverDB.Exec("CREATE DATABASE IF NOT EXISTS `" + name + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
			fmt.Printf("创建数据库 %s 失败: %v\n", name, err)
			return
		}
	}

	migrations := []struct {
		name string
		dsn  string
		run  func(*gorm.DB) error
	}{
		{
			name: "beaver_account",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_account?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&user_models.UserModel{}, &user_models.UserChangeLogModel{}, &user_models.UserSettingModel{},
					&auth_models.AuthCredentialModel{}, &auth_models.AuthDeviceModel{},
				)
			},
		},
		{
			name: "beaver_relation",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_relation?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&friend_models.FriendModel{},
					&friend_models.FriendVerifyModel{},
					&friend_models.FriendBlockModel{},
					&group_models.GroupModel{},
					&group_models.GroupMemberModel{},
					&group_models.GroupJoinRequestModel{},
					&group_models.GroupMemberChangeLogModel{},
					&group_models.GroupInviteLinkModel{},
				)
			},
		},
		{
			name: "beaver_message",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_message?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&message_models.ChatMessage{},
					&message_models.ChatConversationMeta{},
					&message_models.ChatUserConversation{},
					&message_models.ChatUserDelete{},
					&message_models.ChatMessageMedia{},
					&message_models.ChatForward{},
				)
			},
		},
		{
			name: "beaver_asset",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_asset?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&emoji_models.EmojiPackage{},
					&emoji_models.Emoji{},
					&emoji_models.EmojiPackageEmoji{},
					&emoji_models.EmojiPackageCollect{},
					&emoji_models.EmojiCollectEmoji{},
					&file_models.FileModel{},
				)
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
			name: "beaver_platform",
			dsn:  "root:123456@tcp(127.0.0.1:3306)/beaver_platform?charset=utf8mb4&parseTime=True&loc=Local",
			run: func(db *gorm.DB) error {
				return db.AutoMigrate(
					&track_models.TrackBucket{},
					&track_models.TrackLogger{},
					&feedback_models.FeedbackModel{},
					&report_models.ContentReportModel{},
					&update_models.UpdateApp{},
					&update_models.UpdateArchitecture{},
					&update_models.UpdateVersion{},
					&update_models.UpdateReleasePolicy{},
					&update_models.UpdateReport{},
					&workbench_models.WorkbenchApp{},
					&org_models.Organization{},
					&org_models.Department{},
					&org_models.Staff{},
					&org_models.DeptUser{},
					&org_models.Admin{},
					// AI/知识库业务表已随 AI 服务拆分至 beaver-ai（自建 beaver_ai 库）
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

	fileDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_asset?charset=utf8mb4&parseTime=True&loc=Local")
	platformDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_platform?charset=utf8mb4&parseTime=True&loc=Local")
	userDB := coregorm.InitGorm("root:123456@tcp(127.0.0.1:3306)/beaver_account?charset=utf8mb4&parseTime=True&loc=Local")
	authDB := userDB
	_ = fileseed.InitDefaultFiles(fileDB)
	_ = platform.InitPlatform(platformDB)
	if err := userseed.InitDefaultUser(userDB, authDB); err != nil {
		fmt.Printf("默认用户初始化失败: %v\n", err)
		return
	}
	fmt.Println("所有库表结构生成完成")
}
