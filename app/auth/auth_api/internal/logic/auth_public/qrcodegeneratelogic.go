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

package auth_public

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	util "beaver/utils/uuid"
)

const (
	QrcodeStatusPending   = "pending"
	QrcodeStatusConfirmed = "confirmed"
	QrcodeStatusExpired   = "expired"

	QrcodeTTL              = 3 * time.Minute
	QrcodeKeyFmt           = "qrcode:%s"
	QrcodeTokenExpireHours = 12
)

type QrcodeSession struct {
	Status        string `json:"status"`
	Source        string `json:"source"`
	ScannedUserID string `json:"scannedUserId,omitempty"`
}

type QrcodeGenerateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewQrcodeGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeGenerateLogic {
	return &QrcodeGenerateLogic{
		ctx:    ctx,
		logger: beaverlog.New("qrcode_generate", ctx),
		svcCtx: svcCtx,
	}
}

func (l *QrcodeGenerateLogic) QrcodeGenerate(req *types.QrcodeGenerateReq) (*types.QrcodeGenerateRes, error) {
	source := req.Source
	if source == "" {
		source = "login"
	}

	token := util.NewV4().String()
	session := QrcodeSession{
		Status: QrcodeStatusPending,
		Source: source,
	}
	sessionJSON, _ := json.Marshal(session)
	key := fmt.Sprintf(QrcodeKeyFmt, token)

	if err := l.svcCtx.Redis.Set(key, string(sessionJSON), QrcodeTTL).Err(); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "二维码写入失败",
			Data: map[string]any{"err": err.Error()},
		})
		return nil, fmt.Errorf("服务内部异常")
	}

	l.logger.Info(model.LogMsg{
		Text: "二维码生成成功",
		Data: map[string]any{"source": source},
	})

	return &types.QrcodeGenerateRes{
		Token:    token,
		ExpireAt: time.Now().Add(QrcodeTTL).Unix(),
	}, nil
}
