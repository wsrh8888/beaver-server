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
	"time"

	user_models "beaver/app/account/account_models/user"
	"beaver/app/account/account_rpc/internal/svc"
	"beaver/app/account/account_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	utils "beaver/utils/rand"
)

// 工号序号统一从 100000 起；Redis 计数器里存"当前已用到的最大序号"。
const minUserID = 100000

// userIDRule 是某类用户 user_id（即工号）的生成规则：
// prefix 为单字母类型前缀，key 为该类型独立的 Redis 自增计数器。
type userIDRule struct {
	prefix string
	key    string
}

// userIDRules 是工号规则的唯一注册表：单字母类型前缀 + 6 位自增序号。
// 普通用户不加前缀；新增用户类型只需在这里加一行。
var userIDRules = map[int8]userIDRule{
	user_models.UserTypeNormal:  {prefix: "", key: "user_id_counter"},
	user_models.UserTypeBot:     {prefix: "B", key: "user_id_bot_counter"},   // B = Bot，外接机器人
	user_models.UserTypeDigital: {prefix: "A", key: "user_id_agent_counter"}, // A = Agent，数字伙伴（智能体）
}

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

// generateUserID 按用户类型生成 user_id（即工号）：单字母类型前缀 + 6 位自增序号。
// 普通用户不加前缀（100000），Agent 为 A100000，Bot 为 B100000。
func (l *UserCreateLogic) generateUserID(userType int8) (string, error) {
	rule, ok := userIDRules[userType]
	if !ok {
		// 未知类型直接失败，避免静默生成普通用户号造成号段污染。
		return "", fmt.Errorf("未知用户类型: %d", userType)
	}

	// 计数器首次初始化为 minUserID-1，INCR 后从 minUserID 开始；已存在则沿用它。
	if err := l.svcCtx.Redis.SetNX(rule.key, minUserID-1, 0).Err(); err != nil {
		return "", fmt.Errorf("初始化用户ID计数器失败: %v", err)
	}

	result := l.svcCtx.Redis.Incr(rule.key)
	if result.Err() != nil {
		l.logger.Error(model.LogMsg{
			Text: "生成用户ID失败",
			Data: map[string]any{"prefix": rule.prefix, "err": result.Err().Error()},
		})
		return "", fmt.Errorf("生成用户ID失败: %v", result.Err())
	}

	if rule.prefix == "" {
		return strconv.FormatInt(result.Val(), 10), nil
	}
	return fmt.Sprintf("%s%06d", rule.prefix, result.Val()), nil
}

func (l *UserCreateLogic) UserCreate(in *pb.UserCreateReq) (*pb.UserCreateRes, error) {
	// UserType 必传，由 generateUserID 的号段注册表校验；不传直接失败，不做默认兜底。
	userType := int8(in.UserType)

	// 验证必填字段（Bot / Digital 用户不需要联系方式）
	isSystemUser := userType == user_models.UserTypeBot || userType == user_models.UserTypeDigital
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

	// 生成 user_id（即工号）：普通用户 100000、Agent A100000、Bot B100000。
	userID, genErr := l.generateUserID(userType)
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
		UserType: userType,
		Version:  version,
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
