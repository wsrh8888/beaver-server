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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SaveDigitalKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDigitalKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalKnowledgeLogic {
	return &SaveDigitalKnowledgeLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_digital_knowledge", ctx)}
}

// SaveDigitalKnowledge 新建或更新知识库。
// 新建时归属记在当前开发者名下（acl 内部共享），并顺手绑定到该数字伙伴，
// 否则用户建完还得再点一次「绑定」，白填一趟。
func (l *SaveDigitalKnowledgeLogic) SaveDigitalKnowledge(req *types.SaveDigitalKnowledgeReq) (resp *types.SaveDigitalKnowledgeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if req.CollectionID == "" {
		if name == "" {
			return nil, errors.New("知识库名称不能为空")
		}
		created, err := l.svcCtx.PlatformRpc.CreateCollection(l.ctx, &platform.CreateCollectionReq{
			Name:        name,
			Kind:        "docs",
			AclMode:     "internal",
			Description: strings.TrimSpace(req.Description),
			Status:      1,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "创建知识库失败", Data: map[string]interface{}{"err": err.Error()}})
			return nil, errors.New("创建知识库失败")
		}
		collectionID := created.CollectionId
		if req.TopK > 0 {
			if err := l.updateIndexConfig(collectionID, "", req.TopK); err != nil {
				return nil, err
			}
		}
		if err := l.bindKnowledge(partner.DigitalID, collectionID); err != nil {
			return nil, err
		}
		return &types.SaveDigitalKnowledgeRes{CollectionID: collectionID}, nil
	}

	// 更新：先读原配置，indexConfig 只覆盖传入的键，避免把 embedding 凭据写没
	current, err := l.svcCtx.PlatformRpc.GetCollection(l.ctx, &platform.GetCollectionReq{CollectionId: req.CollectionID})
	if err != nil || current.Collection == nil {
		return nil, errors.New("知识库不存在")
	}
	indexConfig := mergeIndexConfig(current.Collection.IndexConfig, req.IndexConfig)
	if req.TopK > 0 {
		indexConfig = setIndexConfigTopK(indexConfig, req.TopK)
	}
	update := &platform.UpdateCollectionReq{
		CollectionId: req.CollectionID,
		IndexConfig:  indexConfig,
	}
	if name != "" {
		update.Name = name
	}
	if req.Description != "" {
		update.Description = req.Description
	}
	if _, err := l.svcCtx.PlatformRpc.UpdateCollection(l.ctx, update); err != nil {
		l.logger.Error(model.LogMsg{Text: "更新知识库失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("更新知识库失败")
	}
	return &types.SaveDigitalKnowledgeRes{CollectionID: req.CollectionID}, nil
}

func (l *SaveDigitalKnowledgeLogic) updateIndexConfig(collectionID, existing string, topK int) error {
	_, err := l.svcCtx.PlatformRpc.UpdateCollection(l.ctx, &platform.UpdateCollectionReq{
		CollectionId: collectionID,
		IndexConfig:  setIndexConfigTopK(existing, topK),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "更新知识库检索条数失败", Data: map[string]interface{}{"err": err.Error()}})
		return errors.New("更新知识库失败")
	}
	return nil
}

// bindKnowledge 把知识库追加到数字伙伴的绑定列表（已在列表里就不动）。
func (l *SaveDigitalKnowledgeLogic) bindKnowledge(digitalID, collectionID string) error {
	ids := loadDigitalKnowledgeIDs(l.svcCtx.DB, digitalID)
	for _, id := range ids {
		if id == collectionID {
			return nil
		}
	}
	ids = append(ids, collectionID)
	return l.writeKnowledgeIDs(digitalID, ids)
}

func (l *SaveDigitalKnowledgeLogic) writeKnowledgeIDs(digitalID string, ids []string) error {
	return writeDigitalKnowledgeIDs(l.svcCtx.DB, digitalID, ids)
}
