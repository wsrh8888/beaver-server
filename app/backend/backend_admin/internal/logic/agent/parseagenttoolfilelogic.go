/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package agent

import (
	"context"
	"encoding/base64"
	"path"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// maxSkillFileBytes SKILL.md 大小上限（2MB），一个技能说明不该比这更大。
const maxSkillFileBytes = 2 * 1024 * 1024

type ParseAgentToolFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewParseAgentToolFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseAgentToolFileLogic {
	return &ParseAgentToolFileLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("ParseAgentToolFile", ctx)}
}

// ParseAgentToolFile 解析上传的 SKILL.md。
// 技能就是一个文件：frontmatter 给 name/description/version，正文原样回给前端预览，
// 确认后前端再调 tool_create 入库 —— 这里不落库，避免传错文件留下脏数据。
func (l *ParseAgentToolFileLogic) ParseAgentToolFile(req *types.ParseAgentToolFileReq) (resp *types.ParseAgentToolFileRes, err error) {
	fileName := strings.TrimSpace(req.FileName)
	if fileName == "" {
		return &types.ParseAgentToolFileRes{Msg: "请选择 SKILL.md 文件"}, nil
	}
	ext := strings.ToLower(path.Ext(fileName))
	if ext != ".md" && ext != ".markdown" && ext != ".txt" {
		return &types.ParseAgentToolFileRes{Msg: "只支持 .md / .markdown / .txt"}, nil
	}

	data, decErr := base64.StdEncoding.DecodeString(req.Data)
	if decErr != nil {
		return &types.ParseAgentToolFileRes{Msg: "文件内容不是合法的 base64"}, nil
	}
	if len(data) == 0 {
		return &types.ParseAgentToolFileRes{Msg: "文件是空的"}, nil
	}
	if len(data) > maxSkillFileBytes {
		return &types.ParseAgentToolFileRes{Msg: "文件超过 2MB 上限"}, nil
	}

	doc, parseErr := parseSkillFile(fileName, string(data))
	if parseErr != nil {
		l.logger.Warn(model.LogMsg{Text: "技能文件解析失败", Data: map[string]interface{}{"file": fileName, "err": parseErr.Error()}})
		return &types.ParseAgentToolFileRes{Msg: parseErr.Error()}, nil
	}

	return &types.ParseAgentToolFileRes{
		Name:        doc.Name,
		Description: doc.Description,
		Version:     doc.Version,
		SkillMD:     doc.Body,
		FileName:    path.Base(fileName),
		Config:      doc.ConfigJSON,
	}, nil
}
