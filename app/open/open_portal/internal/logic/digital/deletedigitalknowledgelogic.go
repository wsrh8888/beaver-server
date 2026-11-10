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
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteDigitalKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalKnowledgeLogic {
	return &DeleteDigitalKnowledgeLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital_knowledge", ctx)}
}

// DeleteDigitalKnowledge 删除知识库（连同其下文档与索引）；传入 digitalId 时顺带从该数字伙伴解绑。
// 知识库是平台级资源，删掉后其它数字伙伴也会失去它——前端必须做二次确认。
func (l *DeleteDigitalKnowledgeLogic) DeleteDigitalKnowledge(req *types.DeleteDigitalKnowledgeReq) (resp *types.DeleteDigitalKnowledgeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.DigitalID != "" {
		if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
			return nil, err
		}
	}
	if req.CollectionID == "" {
		return nil, errors.New("collectionId 不能为空")
	}

	if _, err := l.svcCtx.PlatformRpc.DeleteCollection(l.ctx, &platform.DeleteCollectionReq{
		CollectionId: req.CollectionID,
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识库失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("删除知识库失败")
	}

	if req.DigitalID != "" {
		remain := make([]string, 0)
		for _, id := range loadDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID) {
			if id != req.CollectionID {
				remain = append(remain, id)
			}
		}
		if err := writeDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID, remain); err != nil {
			l.logger.Error(model.LogMsg{Text: "解绑已删除的知识库失败", Data: map[string]interface{}{"err": err.Error()}})
		}
	}
	return &types.DeleteDigitalKnowledgeRes{}, nil
}
