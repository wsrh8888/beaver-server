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

package oauth_public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PasswordLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码登录
func NewPasswordLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordLoginLogic {
	return &PasswordLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordLoginLogic) PasswordLogin(req *types.PasswordLoginReq) (resp *types.PasswordLoginRes, err error) {
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	userRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Account,
		Type:    "email",
	})
	if err != nil {
		userRes, err = l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
			Keyword: req.Account,
			Type:    "userId",
		})
	}
	if err != nil {
		logx.Errorf("用户不存在: account=%s, err=%v", req.Account, err)
		return nil, errors.New("账号或密码错误")
	}

	verifyRes, err := l.svcCtx.AuthRpc.VerifyPassword(l.ctx, &auth_rpc.VerifyPasswordReq{
		UserId:   userRes.UserInfo.UserId,
		Password: req.Password,
	})
	if err != nil || !verifyRes.Valid {
		logx.Errorf("密码错误: userId=%s", userRes.UserInfo.UserId)
		return nil, errors.New("账号或密码错误")
	}

	code, expireIn, err := createOAuthCode(l.svcCtx.DB, req.AppID, userRes.UserInfo.UserId, "password", "")
	if err != nil {
		logx.Errorf("生成授权码失败: userId=%s, err=%v", userRes.UserInfo.UserId, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	logx.Infof("账号密码登录成功: userId=%s, appId=%s, account=%s", userRes.UserInfo.UserId, req.AppID, req.Account)

	return &types.PasswordLoginRes{
		Code:     code,
		ExpireIn: expireIn,
	}, nil
}

func createOAuthCode(db *gorm.DB, appID, userID, scene, sceneRef string) (string, int64, error) {
	if appID == "" || userID == "" {
		return "", 0, errors.New("参数不完整")
	}

	var app open_models.OpenApp
	if err := db.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error; err != nil {
		return "", 0, errors.New("应用不存在或未启用")
	}

	var oauthConfig open_models.OpenAppOAuth
	scope := ""
	if err := db.Where("app_id = ?", appID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scope = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scope = string(data)
	}

	const ttl = 5 * time.Minute
	code := util.NewV4().String()
	record := open_models.OpenOAuthCode{
		Code:      code,
		AppID:     appID,
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Scene:     scene,
		State:     sceneRef,
	}
	if err := db.Create(&record).Error; err != nil {
		return "", 0, errors.New("生成授权码失败")
	}
	return code, int64(ttl.Seconds()), nil
}
