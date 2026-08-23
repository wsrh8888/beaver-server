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

package logic

import (
	"context"
	"encoding/json"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProxySendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProxySendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxySendMsgLogic {
	return &ProxySendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxySendMsgLogic) ProxySendMsg(req *types.ProxySendMsgReq) (resp *types.ProxySendMsgRes, err error) {
	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		return nil, err
	}

	content := type_struct.WsContent{
		Data: type_struct.WsData{
			Type:           wsTypeConst.Type(req.Type),
			Body:           bodyBytes,
			ConversationID: req.ConversationId,
		},
	}

	ws_conn.SendMsgToUser(req.TargetID, wsCommandConst.Command(req.Command), content)

	return &types.ProxySendMsgRes{}, nil
}
