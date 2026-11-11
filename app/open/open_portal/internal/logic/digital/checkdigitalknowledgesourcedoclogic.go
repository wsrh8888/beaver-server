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

	"beaver/app/doc/doc_rpc/doc"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CheckDigitalKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCheckDigitalKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckDigitalKnowledgeSourceDocLogic {
	return &CheckDigitalKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("check_digital_knowledge_source_doc", ctx)}
}

// CheckDigitalKnowledgeSourceDoc 检查用户粘贴的在线文档地址：解析出 doc_id 后
// 按当前用户调 doc_rpc.GetDoc 鉴权，成功则回带文档标题供前端在弹窗里预览。
// 只校验、不建立任何引用；真正的导入由 doc/import 完成。
func (l *CheckDigitalKnowledgeSourceDocLogic) CheckDigitalKnowledgeSourceDoc(req *types.CheckDigitalKnowledgeSourceDocReq) (resp *types.CheckDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	docID := parseBeaverDocID(req.DocURL)
	if docID == "" {
		return &types.CheckDigitalKnowledgeSourceDocRes{
			Valid: false,
			Msg:   "无法从链接中解析出文档 ID，请粘贴海狸在线文档地址或裸文档 ID",
		}, nil
	}

	res, authErr := l.svcCtx.DocRpc.GetDoc(l.ctx, &doc.GetDocReq{
		DocId:  docID,
		UserId: req.UserID,
	})
	if authErr != nil || res == nil || res.Doc == nil {
		l.logger.Warn(model.LogMsg{Text: "在线文档检查失败", Data: map[string]interface{}{"docId": docID, "err": authErr}})
		return &types.CheckDigitalKnowledgeSourceDocRes{
			Valid: false,
			DocID: docID,
			Msg:   "文档不存在，或你没有查看权限",
		}, nil
	}

	title := strings.TrimSpace(res.Doc.Title)
	if title == "" {
		title = docID
	}
	return &types.CheckDigitalKnowledgeSourceDocRes{
		Valid: true,
		DocID: docID,
		Title: title,
	}, nil
}
