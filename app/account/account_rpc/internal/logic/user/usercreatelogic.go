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

package userlogic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	user_models "beaver/app/account/account_models/user"
	"beaver/app/account/account_rpc/internal/svc"
	"beaver/app/account/account_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	utils "beaver/utils/rand"
)

const (
	minUserID = 100000

	// user_id 即工号，按用户类型加前缀：普通用户纯 6 位数字，bot BOT+6位，数字伙伴 DG+6位。
	// 新增用户类型时在这里加一行前缀 + 一个计数器 key 即可。
	userIDKeyNormal  = "user_id_counter"
	userIDKeyBot     = "user_id_bot_counter"
	userIDKeyDigital = "user_id_digital_counter"
	prefixBot        = "BOT"
	prefixDigital    = "DG"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("user_create", ctx),
	}
}

// generateUserID 按用户类型生成 user_id（即工号）：
// 普通用户纯 6 位数字；bot BOT+6位；数字伙伴 DG+6位。前缀在 const 里集中定义，新增类型加一行即可。
func (l *UserCreateLogic) generateUserID(userType int8) (string, error) {
	key, prefix := userIDKeyNormal, ""
	switch userType {
	case user_models.UserTypeBot:
		key, prefix = userIDKeyBot, prefixBot
	case user_models.UserTypeDigital:
		key, prefix = userIDKeyDigital, prefixDigital
	}

	if err := l.initializeIDCounter(key, prefix); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "初始化用户ID计数器失败",
			Data: map[string]any{"prefix": prefix, "err": err.Error()},
		})
		return "", fmt.Errorf("初始化用户ID计数器失败: %v", err)
	}

	result := l.svcCtx.Redis.Incr(key)
	if result.Err() != nil {
		l.logger.Error(model.LogMsg{
			Text: "生成用户ID失败",
			Data: map[string]any{"prefix": prefix, "err": result.Err().Error()},
		})
		return "", fmt.Errorf("生成用户ID失败: %v", result.Err())
	}

	if prefix == "" {
		return strconv.FormatInt(result.Val(), 10), nil
	}
	return fmt.Sprintf("%s%06d", prefix, result.Val()), nil
}

// initializeIDCounter 从库里同类 user_id 的最大序号续接，保证不与存量撞号。
func (l *UserCreateLogic) initializeIDCounter(key, prefix string) error {
	exists := l.svcCtx.Redis.Exists(key)
	if exists.Err() != nil {
		return fmt.Errorf("检查用户ID计数器失败: %v", exists.Err())
	}
	if exists.Val() != 0 {
		return nil
	}

	var maxUser user_models.UserModel
	var current int64
	if prefix == "" {
		// 普通用户：纯数字 user_id，从 100000 起。
		err := l.svcCtx.DB.Select("user_id").
			Where("user_id REGEXP '^[0-9]+$'").
			Order("CAST(user_id AS UNSIGNED) DESC").
			First(&maxUser).Error
		current = minUserID - 1
		if err == nil && maxUser.UserID != "" {
			if n, parseErr := strconv.ParseInt(maxUser.UserID, 10, 64); parseErr == nil {
				current = n
			}
		} else if err != nil && err.Error() != "record not found" {
			return fmt.Errorf("查询最大用户ID失败: %v", err)
		}
	} else {
		// 带前缀：取同前缀里最大的序号续接。
		start := len(prefix) + 1
		err := l.svcCtx.DB.Select("user_id").
			Where("user_id LIKE ?", prefix+"%").
			Order(fmt.Sprintf("CAST(SUBSTRING(user_id, %d) AS UNSIGNED) DESC", start)).
			First(&maxUser).Error
		if err == nil && maxUser.UserID != "" {
			if n, parseErr := strconv.ParseInt(strings.TrimPrefix(maxUser.UserID, prefix), 10, 64); parseErr == nil {
				current = n
			}
		} else if err != nil && err.Error() != "record not found" {
			return fmt.Errorf("查询最大用户ID失败: %v", err)
		}
	}

	// 计数器存"当前已用到的最大序号"，下一次 INCR 正好得到 +1。
	if err := l.svcCtx.Redis.Set(key, current, 0).Err(); err != nil {
		return fmt.Errorf("设置用户ID计数器失败: %v", err)
	}
	l.logger.Info(model.LogMsg{
		Text: "用户ID计数器初始化完成",
		Data: map[string]interface{}{"prefix": prefix, "current": current},
	})
	return nil
}

func (l *UserCreateLogic) UserCreate(in *pb.UserCreateReq) (*pb.UserCreateRes, error) {
	// 验证必填字段（Bot / Digital 用户不需要联系方式）
	isSystemUser := in.UserType == int32(user_models.UserTypeBot) || in.UserType == int32(user_models.UserTypeDigital)
	if !isSystemUser {
		if in.Phone == "" && in.Email == "" {
			return nil, errors.New("手机号或邮箱至少需要提供一个")
		}
	} else {
		// 机器人用户：昵称必填
		if in.NickName == "" {
			return nil, errors.New("机器人用户必须提供昵称")
		}
	}

	// 检查用户是否已存在（通过 phone/email 查询）
	var user user_models.UserModel
	var err error

	if in.Phone != "" && in.Email != "" {
		// 手机号和邮箱都提供，检查是否已存在
		err = l.svcCtx.DB.Take(&user, "phone = ? OR email = ?", in.Phone, in.Email).Error
	} else if in.Phone != "" {
		// 只提供手机号
		err = l.svcCtx.DB.Take(&user, "phone = ?", in.Phone).Error
	} else if in.Email != "" {
		// 只提供邮箱
		err = l.svcCtx.DB.Take(&user, "email = ?", in.Email).Error
	} else {
		// 机器人用户：没有联系方式，跳过重复检查
		err = errors.New("record not found")
	}

	if err == nil {
		return nil, errors.New("用户已存在")
	}

	// 生成随机昵称
	nickName := utils.GenerateRandomString(8)
	if in.NickName != "" {
		nickName = in.NickName
	}

	// 生成 user_id（即工号）：普通用户 6 位数字、bot BOT+6位、数字伙伴 DG+6位。
	userID, genErr := l.generateUserID(int8(in.UserType))
	if genErr != nil {
		return nil, errors.New("生成用户ID失败")
	}

	// 获取新版本号（用户独立递增，从1开始）
	version := l.svcCtx.VersionGen.GetNextVersion("users", "user_id", userID)
	if version == -1 {
		l.logger.Error(model.LogMsg{
			Text: "获取用户版本号失败",
			Data: map[string]any{"userId": userID},
		})
		return nil, errors.New("获取版本号失败")
	}
	l.logger.Info(model.LogMsg{
		Text: "获取用户版本号",
		Data: map[string]interface{}{"userId": userID, "version": version},
	})

	// 创建用户基础信息（不包含密码）
	user = user_models.UserModel{
		UserID:   userID,
		Email:    in.Email,
		Phone:    in.Phone,
		Source:   in.Source,
		NickName: nickName,
		Abstract: in.Abstract,
		UserType: int8(in.UserType),
		Version:  version,
	}
	if user.UserType == 0 {
		user.UserType = user_models.UserTypeNormal
	}

	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建用户失败",
			Data: map[string]any{"userId": userID, "err": err.Error()},
		})
		return nil, errors.New("创建用户失败")
	}

	// TODO: 调用 auth 服务创建用户凭证（密码）
	// 这里需要通过 RPC 调用 auth 服务的 CreateCredential 接口
	// 暂时留空，等待 auth 服务实现后补充

	// 记录用户创建的变更日志
	l.recordUserCreateLog(user.UserID, version)

	l.logger.Info(model.LogMsg{
		Text: "用户创建成功",
		Data: map[string]interface{}{"userId": user.UserID, "version": version},
	})

	return &pb.UserCreateRes{
		UserID: user.UserID,
	}, nil
}

// recordUserCreateLog 记录用户创建日志
func (l *UserCreateLogic) recordUserCreateLog(userID string, version int64) {
	// 创建用户创建的变更日志
	changeLog := user_models.UserChangeLogModel{
		UserID:     userID,
		ChangeType: "create",
		NewValue:   "",
		ChangeTime: time.Now().Unix(),
		Version:    version,
	}

	if err := l.svcCtx.DB.Create(&changeLog).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "记录用户创建日志失败",
			Data: map[string]any{"userId": userID, "err": err.Error()},
		})
	} else {
		l.logger.Info(model.LogMsg{
			Text: "用户创建日志记录成功",
			Data: map[string]interface{}{"userId": userID, "version": version},
		})
	}
}
