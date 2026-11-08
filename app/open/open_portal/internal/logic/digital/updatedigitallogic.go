/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"
	"strings"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type UpdateDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDigitalLogic {
	return &UpdateDigitalLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("update_digital", ctx)}
}

func (l *UpdateDigitalLogic) UpdateDigital(req *types.UpdateDigitalReq) (resp *types.UpdateDigitalRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}
	res := &types.UpdateDigitalRes{}

	if req.Name != "" || req.Avatar != "" {
		_, err = l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId:   partner.DigitalID,
			NickName: req.Name,
			Avatar:   req.Avatar,
		})
		if err != nil {
			return nil, errors.New("更新展示信息失败")
		}
	}
	if req.Description != "" {
		abstract := req.Description
		_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
			UserIds:       []string{partner.DigitalID},
			PatchAbstract: &abstract,
		})
		if err != nil {
			return nil, errors.New("更新简介失败")
		}
	}

	// 主表字段一次性落库：工号、可见范围、编排方式（编排方式随时可切换，旧配置保留不删）
	updates := map[string]any{}
	workID := strings.TrimSpace(req.WorkID)
	if req.ApplyWorkID && partner.WorkID == "" {
		generated, err := applyWorkID(l.svcCtx.DB)
		if err != nil {
			return nil, err
		}
		updates["work_id"] = generated
	} else if workID != "" && workID != partner.WorkID {
		if len(workID) > 32 {
			return nil, errors.New("工号最长 32 个字符")
		}
		if err := checkWorkIDAvailable(l.svcCtx.DB, workID, partner.DigitalID); err != nil {
			return nil, err
		}
		updates["work_id"] = workID
	}
	if req.Visibility != "" {
		updates["visibility"] = req.Visibility
	}
	if req.Runtime != "" && req.Runtime != partner.Runtime {
		updates["runtime"] = req.Runtime
		// builtin 创建的数字伙伴没有长期凭据，切到 external 时补发一把，避免切过去连不上
		if req.Runtime == open_models.DigitalRuntimeExternal && partner.AppKey == "" {
			appKey, appSecret := generateExternalCredentials()
			updates["app_key"] = appKey
			updates["app_secret_hash"] = appSecret
			res.AppKey = appKey
			res.AppSecret = appSecret
		}
	}
	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(partner).Updates(updates).Error; err != nil {
			return nil, errors.New("更新失败")
		}
	}
	return res, nil
}
