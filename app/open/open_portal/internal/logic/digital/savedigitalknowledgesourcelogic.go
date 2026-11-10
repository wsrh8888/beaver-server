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

type SaveDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalKnowledgeSourceLogic {
	return &SaveDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_digital_knowledge_source", ctx)}
}

// SaveDigitalKnowledgeSource 新建知识源，或改名 / 改检索设置。
// 只接受「部分字段」语义：空字符串与 0 都代表不改，避免其它分区的整份提交互相覆盖。
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

	var row open_models.OpenKnowledgeSource
	if req.KnowledgeSourceID == "" {
		if name == "" {
			return nil, errors.New("知识源名称不能为空")
		}
		if topK <= 0 {
			topK = open_models.DefaultKnowledgeSourceTopK
		}
		row = open_models.OpenKnowledgeSource{
			SourceID:      newKnowledgeSourceID(),
			DigitalID:     req.DigitalID,
			OwnerID:       req.UserID,
			Name:          name,
			RetrievalTopK: topK,
			IndexConfig:   strings.TrimSpace(req.IndexConfig),
			KnowledgeIDs:  open_models.KnowledgeIDList{},
		}
	} else {
		loaded, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
		if err != nil {
			return nil, err
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
