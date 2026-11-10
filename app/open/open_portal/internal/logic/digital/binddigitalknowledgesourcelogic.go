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

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type BindDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewBindDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindDigitalKnowledgeSourceLogic {
	return &BindDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("bind_digital_knowledge_source", ctx)}
}

// BindDigitalKnowledgeSource 给一个知识源引入 N 个平台知识库。
// 同一个知识库可以被不同知识源引入，但同一个知识源内部去重；
// 传入的 ID 若已不存在于平台（库被删了）直接忽略，不写入失效引用。
func (l *BindDigitalKnowledgeSourceLogic) BindDigitalKnowledgeSource(req *types.BindDigitalKnowledgeSourceReq) (resp *types.BindDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	if len(req.KnowledgeIDs) == 0 {
		return &types.BindDigitalKnowledgeSourceRes{KnowledgeIDs: append([]string{}, row.KnowledgeIDs...)}, nil
	}

	existing := collectionNameMap(l.ctx, l.svcCtx)
	ids := uniqTrimmed(row.KnowledgeIDs)
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		seen[id] = true
	}
	for _, id := range uniqTrimmed(req.KnowledgeIDs) {
		if _, ok := existing[id]; !ok {
			continue
		}
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	row.KnowledgeIDs = open_models.KnowledgeIDList(ids)
	if err := l.svcCtx.DB.Save(row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "引入知识库失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("引入知识库失败")
	}
	if err := syncDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID); err != nil {
		return nil, errors.New("同步检索配置失败")
	}
	return &types.BindDigitalKnowledgeSourceRes{KnowledgeIDs: ids}, nil
}
