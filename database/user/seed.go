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

package userseed

import (
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"
	"fmt"
	"log"

	"gorm.io/gorm"
)

const (
	defaultUserPassword = "e10adc3949ba59abbe56e057f20f883e" // MD5("123456")
	seedEmailStart      = 751135380
	seedEmailEnd        = 751135389
	ownerEmailLocal     = 751135385
	ownerUserID         = "100000"
)

// seedAccounts 返回初始化账号列表（邮箱本地段 → userId）
// 751135385 固定为 100000（开放平台默认 owner），其余顺延 100001~100009
func seedAccounts() []struct {
	email  string
	userID string
	nick   string
} {
	accounts := make([]struct {
		email  string
		userID string
		nick   string
	}, 0, seedEmailEnd-seedEmailStart+1)

	nextID := 100001
	for local := seedEmailStart; local <= seedEmailEnd; local++ {
		userID := fmt.Sprintf("%d", nextID)
		if local == ownerEmailLocal {
			userID = ownerUserID
		} else {
			nextID++
		}
		accounts = append(accounts, struct {
			email  string
			userID string
			nick   string
		}{
			email:  fmt.Sprintf("%d@qq.com", local),
			userID: userID,
			nick:   fmt.Sprintf("测试用户%d", local%1000),
		})
	}
	return accounts
}

// InitDefaultUser 初始化默认测试账号（用户 + 密码）
func InitDefaultUser(userDB, authDB *gorm.DB) error {
	for _, acc := range seedAccounts() {
		if err := ensureUser(userDB, acc.email, acc.userID, acc.nick); err != nil {
			return err
		}
		if err := ensureCredential(authDB, acc.userID); err != nil {
			return err
		}
	}
	return nil
}

func ensureUser(userDB *gorm.DB, email, userID, nickName string) error {
	var user user_models.UserModel
	err := userDB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = user_models.UserModel{
			UserID:   userID,
			UserType: user_models.UserTypeNormal,
			NickName: nickName,
			Email:    email,
			Source:   user_models.SourceEmail,
			Status:   1,
			Version:  1,
		}
		if err := userDB.Create(&user).Error; err != nil {
			return fmt.Errorf("创建默认用户失败(%s): %w", email, err)
		}
		log.Printf("创建默认用户成功: userId=%s, email=%s", user.UserID, user.Email)
		return nil
	}
	if err != nil {
		return fmt.Errorf("查询默认用户失败(%s): %w", email, err)
	}
	log.Printf("默认用户已存在: userId=%s, email=%s", user.UserID, user.Email)
	return nil
}

func ensureCredential(authDB *gorm.DB, userID string) error {
	var credential auth_models.AuthCredentialModel
	err := authDB.Where("user_id = ?", userID).First(&credential).Error
	if err == gorm.ErrRecordNotFound {
		credential = auth_models.AuthCredentialModel{
			UserID:   userID,
			Password: pwd.HahPwd(defaultUserPassword),
		}
		if err := authDB.Create(&credential).Error; err != nil {
			return fmt.Errorf("创建默认用户凭证失败(userId=%s): %w", userID, err)
		}
		log.Printf("创建默认用户凭证成功: userId=%s", userID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("查询默认用户凭证失败(userId=%s): %w", userID, err)
	}
	log.Printf("默认用户凭证已存在: userId=%s", userID)
	return nil
}
