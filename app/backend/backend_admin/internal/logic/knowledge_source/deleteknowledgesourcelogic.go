/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package knowledge_source

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeSourceLogic {
	return &DeleteKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("DeleteKnowledgeSource", ctx)}
}

// DeleteKnowledgeSource 删除知识源本身与它持有的索引集合，并回写运行时检索列表。
// 集合删除失败不阻断行删除：集合已不存在时也要允许把这条脏知识源清掉。
func (l *DeleteKnowledgeSourceLogic) DeleteKnowledgeSource(req *types.DeleteKnowledgeSourceReq) (resp *types.DeleteKnowledgeSourceRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}
	row, err := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	if collectionID := sourceCollectionID(row.KnowledgeIDs); collectionID != "" {
		if _, rpcErr := l.svcCtx.PlatformRpc.DeleteCollection(l.ctx, &platform.DeleteCollectionReq{
			CollectionId: collectionID,
		}); rpcErr != nil {
			l.logger.Warn(model.LogMsg{Text: "删除知识源索引集合失败，继续删除知识源", Data: map[string]interface{}{"collectionId": collectionID, "err": rpcErr.Error()}})
		}
	}
	if delErr := l.svcCtx.DB.Delete(&row).Error; delErr != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识源失败", Data: map[string]interface{}{"err": delErr.Error(), "sourceId": row.SourceID}})
		return nil, delErr
	}
	if err := syncSceneKnowledgeIDs(l.ctx, l.svcCtx, l.logger, scene); err != nil {
		return nil, err
	}
	return &types.DeleteKnowledgeSourceRes{}, nil
}
