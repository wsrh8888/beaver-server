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
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	doccli "beaver/app/doc/doc_rpc/doc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CheckKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCheckKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckKnowledgeSourceDocLogic {
	return &CheckKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("CheckKnowledgeSourceDoc", ctx)}
}

// CheckKnowledgeSourceDoc 检查运营粘贴的在线文档地址：解析出 doc_id 后按内部接口读取快照，
// 能读到就视为可用。后台是平台管理员视角，不做「某个普通用户是否可见」的校验。
// 只校验、不建立任何引用；真正的导入由 source_doc_import 完成。
func (l *CheckKnowledgeSourceDocLogic) CheckKnowledgeSourceDoc(req *types.CheckKnowledgeSourceDocReq) (resp *types.CheckKnowledgeSourceDocRes, err error) {
	if _, err := requireScene(req.Scene); err != nil {
		return nil, err
	}

	docID := parseBeaverDocID(req.DocURL)
	if docID == "" {
		return &types.CheckKnowledgeSourceDocRes{
			Valid: false,
			Msg:   "无法从链接中解析出文档 ID，请粘贴海狸在线文档地址或裸文档 ID",
		}, nil
	}

	res, readErr := l.svcCtx.DocRpc.GetDocContent(l.ctx, &doccli.GetDocContentReq{DocId: docID})
	if readErr != nil || res == nil {
		l.logger.Warn(model.LogMsg{Text: "在线文档检查失败", Data: map[string]interface{}{"docId": docID, "err": readErr}})
		return &types.CheckKnowledgeSourceDocRes{
			Valid: false,
			DocID: docID,
			Msg:   "文档不存在，或正文尚未生成",
		}, nil
	}

	title := strings.TrimSpace(res.Title)
	if title == "" {
		title = docID
	}
	return &types.CheckKnowledgeSourceDocRes{Valid: true, DocID: docID, Title: title}, nil
}
