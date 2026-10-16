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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetDigitalListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDigitalListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalListLogic {
	return &GetDigitalListLogic{logger: beaverlog.New("get_digital_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDigitalListLogic) GetDigitalList(req *types.GetDigitalListReq) (resp *types.GetDigitalListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListDigitals(l.ctx, &open_rpc.ListDigitalsReq{
		Source:   req.Source,
		Status:   int32(req.Status),
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取数字伙伴列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	list := make([]types.DigitalAdminInfo, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.DigitalAdminInfo{
			DigitalID:   item.DigitalId,
			Name:        item.Name,
			Description: item.Description,
			Avatar:      item.Avatar,
			OwnerID:     item.OwnerId,
			Runtime:     item.Runtime,
			Visibility:  item.Visibility,
			Status:      int(item.Status),
			Source:      item.Source,
			AppID:        item.AppId,
			SystemPrompt: item.SystemPrompt,
			LLMModelID:   item.LlmModelId,
			Tools:        item.Tools,
			McpServerIDs: item.McpServerIds,
			CreatedAt:    item.CreatedAt,
		})
	}
	return &types.GetDigitalListRes{Total: rpcRes.Total, List: list}, nil
}
