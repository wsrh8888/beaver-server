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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UnbindDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUnbindDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindDigitalKnowledgeSourceLogic {
	return &UnbindDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("unbind_digital_knowledge_source", ctx)}
}

// UnbindDigitalKnowledgeSource 移除知识源里的一个知识库引用。
// 只删引用，不删知识库本体；若其它知识源还引用着同一个库，运行时仍会检索它。
func (l *UnbindDigitalKnowledgeSourceLogic) UnbindDigitalKnowledgeSource(req *types.UnbindDigitalKnowledgeSourceReq) (resp *types.UnbindDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	target := strings.TrimSpace(req.CollectionID)
	ids := make([]string, 0, len(row.KnowledgeIDs))
	for _, id := range uniqTrimmed(row.KnowledgeIDs) {
		if id != target {
			ids = append(ids, id)
		}
	}
	row.KnowledgeIDs = open_models.KnowledgeIDList(ids)
	if err := l.svcCtx.DB.Save(row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "移除知识库引用失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("移除引用失败")
	}
	if err := syncDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID); err != nil {
		return nil, errors.New("同步检索配置失败")
	}
	return &types.UnbindDigitalKnowledgeSourceRes{KnowledgeIDs: ids}, nil
}
