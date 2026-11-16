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

package database

import (
	auth_models "beaver/app/account/account_models/auth"
	message_models "beaver/app/message/message_models/chat"
	datasync_models "beaver/app/datasync/datasync_models/datasync"
	emoji_models "beaver/app/asset/asset_models/emoji"
	file_models "beaver/app/asset/asset_models/file"
	"beaver/app/relation/relation_models/friend"
	group_models "beaver/app/relation/relation_models/group"
	notification_models "beaver/app/notification/notification_models/notification"
	"beaver/app/platform/platform_models/feedback"
	"beaver/app/platform/platform_models/track"
	"beaver/app/platform/platform_models/update"
	"beaver/app/platform/platform_models/workbench"
	org_models "beaver/app/platform/platform_models/org"
	user_models "beaver/app/account/account_models/user"
	"beaver/core/coregorm"
	fileseed "beaver/database/file"
	platformseed "beaver/database/platform"
	"fmt"

	"gorm.io/gorm"
)

const (
	mysqlUser = "root:123456"
	mysqlAddr = "127.0.0.1:3306"
)

// Migration 单库迁移配置：库名 + 表模型 + 可选默认数据初始化
type Migration struct {
	Name   string
	Models []any
	Init   func(*gorm.DB) error
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
			Name: "beaver_account",
			Models: []any{
				&user_models.UserModel{},
				&user_models.UserChangeLogModel{},
				&user_models.UserSettingModel{},
				&auth_models.AuthCredentialModel{},
				&auth_models.AuthDeviceModel{},
			},
		},
		{
			Name: "beaver_relation",
			Models: []any{
				&friend_models.FriendModel{},
				&friend_models.FriendVerifyModel{},
				&friend_models.FriendBlockModel{},
				&group_models.GroupModel{},
				&group_models.GroupMemberModel{},
				&group_models.GroupJoinRequestModel{},
				&group_models.GroupMemberChangeLogModel{},
			},
		},
		{
			Name: "beaver_message",
			Models: []any{
				&message_models.ChatMessage{},
				&message_models.ChatConversationMeta{},
				&message_models.ChatUserConversation{},
				&message_models.ChatUserDelete{},
				&message_models.ChatMessageMedia{},
				&message_models.ChatForward{},
			},
		},
		{
			Name: "beaver_asset",
			Models: []any{
				&emoji_models.EmojiPackage{},
				&emoji_models.Emoji{},
				&emoji_models.EmojiPackageEmoji{},
				&emoji_models.EmojiPackageCollect{},
				&emoji_models.EmojiCollectEmoji{},
				&file_models.FileModel{},
			},
			Init: fileseed.InitDefaultFiles,
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
			Name: "beaver_platform",
			Models: []any{
				&track_models.TrackBucket{},
				&track_models.TrackLogger{},
				&feedback_models.FeedbackModel{},
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
			},
			Init: platformseed.InitPlatform,
		},
		{
			Name:   "beaver_datasync",
			Models: []any{&datasync_models.DatasyncModel{}},
		},
	}
}

// RunMigrations 建库、迁表、单库默认数据初始化，返回各库连接供跨库初始化使用
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

		if m.Init != nil {
			if err := m.Init(db); err != nil {
				return nil, fmt.Errorf("%s 默认数据初始化失败: %w", m.Name, err)
			}
		}
	}

	return dbMap, nil
}
