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
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SaveDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalKnowledgeSourceLogic {
	return &SaveDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_digital_knowledge_source", ctx)}
}

// SaveDigitalKnowledgeSource 新建知识源（自动建一个 platform 索引集合），或改名 / 改检索设置。
// 只接受「部分字段」语义：空字符串与 0 都代表不改，避免其它分区的整份提交互相覆盖。
//
// 知识源 ≠ 文档本体：文档正文的权威源在 doc 域在线文档库（普通用户所有），新建时只建一个
// platform 索引集合，用户后续粘贴的在线文档地址会作为引用挂在这个集合上。
func (l *SaveDigitalKnowledgeSourceLogic) SaveDigitalKnowledgeSource(req *types.SaveDigitalKnowledgeSourceReq) (resp *types.SaveDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	topK := req.RetrievalTopK
	if topK < 0 || topK > open_models.MaxKnowledgeSourceTopK {
		return nil, errors.New("检索条数超出范围")
	}
	// 前端创建知识源时只传名称，配置留空。这里统一落 {} 占位：
	// platform 的 index_config 是 json 列（空串非法），本表是 text，
	// 两处保持一致，详情接口回显给向导时也能直接当合法 JSON 解析。
	indexConfig := strings.TrimSpace(req.IndexConfig)
	if indexConfig == "" {
		indexConfig = "{}"
	}

	var row open_models.OpenKnowledgeSource
	if req.KnowledgeSourceID == "" {
		if name == "" {
			return nil, errors.New("知识源名称不能为空")
		}
		if topK <= 0 {
			topK = open_models.DefaultKnowledgeSourceTopK
		}
		// 新建知识源时同步建一个 platform 索引集合，文档引用（在线文档）都挂在这个集合上。
		collRes, collErr := l.svcCtx.PlatformRpc.CreateCollection(l.ctx, &platform.CreateCollectionReq{
			Name:        name,
			Kind:        "open_knowledge_source",
			AclMode:     "internal",
			IndexConfig: indexConfig,
		})
		if collErr != nil {
			l.logger.Error(model.LogMsg{Text: "创建知识源索引集合失败", Data: map[string]interface{}{"err": collErr.Error()}})
			return nil, errors.New("创建知识源失败")
		}
		row = open_models.OpenKnowledgeSource{
			SourceID:      newKnowledgeSourceID(),
			DigitalID:     req.DigitalID,
			OwnerID:       req.UserID,
			Name:          name,
			RetrievalTopK: topK,
			IndexConfig:   indexConfig,
			KnowledgeIDs:  open_models.KnowledgeIDList{collRes.CollectionId},
		}
	} else {
		loaded, loadErr := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
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
		l.logger.Error(model.LogMsg{Text: "保存知识源失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("保存知识源失败")
	}
	// 新建时没有引用，改设置不影响引用；统一走一次同步，保证运行时列表与表一致。
	if err := syncDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID); err != nil {
		return nil, errors.New("同步检索配置失败")
	}
	return &types.SaveDigitalKnowledgeSourceRes{KnowledgeSourceID: row.SourceID}, nil
}
