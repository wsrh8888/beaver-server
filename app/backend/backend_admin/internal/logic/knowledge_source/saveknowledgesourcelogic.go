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
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SaveKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveKnowledgeSourceLogic {
	return &SaveKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("SaveKnowledgeSource", ctx)}
}

// SaveKnowledgeSource 新建知识源（同时建一个 platform 索引集合），或改名 / 改检索设置。
// 只接受「部分字段」语义：空字符串与 0 都代表不改，避免其它分区的整份提交互相覆盖。
func (l *SaveKnowledgeSourceLogic) SaveKnowledgeSource(req *types.SaveKnowledgeSourceReq) (resp *types.SaveKnowledgeSourceRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	topK := req.RetrievalTopK
	if topK < 0 || topK > backend_models.MaxAgentKnowledgeSourceTopK {
		return nil, errors.New("检索条数超出范围")
	}

	var row backend_models.AgentKnowledgeSource
	if req.KnowledgeSourceID == "" {
		if name == "" {
			return nil, errors.New("知识源名称不能为空")
		}
		if topK <= 0 {
			topK = backend_models.DefaultAgentKnowledgeSourceTopK
		}
		// 新建知识源时同步建一个 platform 索引集合，文档引用（在线文档 / 本地文件）都挂在这个集合上。
		collRes, collErr := l.svcCtx.PlatformRpc.CreateCollection(l.ctx, &platform.CreateCollectionReq{
			Name:        name,
			Kind:        "agent_knowledge_source",
			AclMode:     "internal",
			IndexConfig: strings.TrimSpace(req.IndexConfig),
		})
		if collErr != nil {
			l.logger.Error(model.LogMsg{Text: "创建知识源索引集合失败", Data: map[string]interface{}{"err": collErr.Error(), "scene": scene}})
			return nil, errors.New("创建知识源失败")
		}
		row = backend_models.AgentKnowledgeSource{
			SourceID:      newKnowledgeSourceID(),
			Scene:         scene,
			Name:          name,
			RetrievalTopK: topK,
			IndexConfig:   strings.TrimSpace(req.IndexConfig),
			KnowledgeIDs:  backend_models.KnowledgeIDList{collRes.CollectionId},
		}
	} else {
		loaded, loadErr := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID)
		if loadErr != nil {
			return nil, loadErr
		}
		row = *loaded
		if name != "" {
			row.Name = name
		}
		if topK > 0 {
			row.RetrievalTopK = topK
		}
		if cfg := strings.TrimSpace(req.IndexConfig); cfg != "" {
			row.IndexConfig = cfg
		}
	}

	if err := l.svcCtx.DB.Save(&row).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "保存知识源失败", Data: map[string]interface{}{"err": err.Error(), "scene": scene}})
		return nil, errors.New("保存知识源失败")
	}
	// 新建时没有引用，改设置不影响引用；统一走一次同步，保证运行时列表与表一致。
	if err := syncSceneKnowledgeIDs(l.ctx, l.svcCtx, l.logger, scene); err != nil {
		return nil, err
	}
	return &types.SaveKnowledgeSourceRes{KnowledgeSourceID: row.SourceID}, nil
}
