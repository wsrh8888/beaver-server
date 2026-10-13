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

package digital

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/user"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDigitalConfigLogic {
	return &UpdateDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("update_digital_config", ctx),
	}
}

func (l *UpdateDigitalConfigLogic) UpdateDigitalConfig(req *types.UpdateDigitalConfigReq) (resp *types.UpdateDigitalConfigRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	partner, err := ensurePortalAppExternalPartner(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app)
	if err != nil {
		return nil, errors.New("更新数字伙伴配置失败")
	}

	if req.Status != nil {
		st := int8(*req.Status)
		if st != 1 && st != 2 {
			return nil, errors.New("status 仅支持 1启用 / 2停用")
		}
		if err := l.svcCtx.DB.Model(partner).Update("status", st).Error; err != nil {
			return nil, errors.New("更新数字伙伴状态失败")
		}
	}

	if req.DigitalName != "" || req.Avatar != "" {
		displayReq := &user.UserUpdateDisplayReq{UserId: partner.PartnerID}
		if req.DigitalName != "" {
			displayReq.NickName = req.DigitalName
		}
		if req.Avatar != "" {
			displayReq.Avatar = req.Avatar
		}
		if _, err := l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, displayReq); err != nil {
			l.logger.Error(model.LogMsg{
				Text: "同步数字伙伴 IM 展示信息失败",
				Data: map[string]interface{}{"partner": partner.PartnerID, "err": err.Error()},
			})
			return nil, errors.New("配置已保存，但同步 IM 昵称/头像失败")
		}
	}

	return &types.UpdateDigitalConfigRes{}, nil
}
