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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalKnowledgeSourceLogic {
	return &DeleteDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital_knowledge_source", ctx)}
}

// DeleteDigitalKnowledgeSource 只删知识源这一行和它派生的运行时引用，
// 被引用的平台知识库本体（集合 / 文档 / 索引）不受影响。
func (l *DeleteDigitalKnowledgeSourceLogic) DeleteDigitalKnowledgeSource(req *types.DeleteDigitalKnowledgeSourceReq) (resp *types.DeleteDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DB.Delete(row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识源失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("删除知识源失败")
	}
	if err := syncDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID); err != nil {
		return nil, errors.New("同步检索配置失败")
	}
	return &types.DeleteDigitalKnowledgeSourceRes{}, nil
}
