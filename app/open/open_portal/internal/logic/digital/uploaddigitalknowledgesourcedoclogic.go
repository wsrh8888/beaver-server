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
	"encoding/base64"
	"errors"
	"path"
	"strings"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// maxKnowledgeUploadBytes 单个本地文件的 base64 解码后上限（20MB，与 doc 域导入一致）。
const maxKnowledgeUploadBytes = 20 * 1024 * 1024

type UploadDigitalKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUploadDigitalKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDigitalKnowledgeSourceDocLogic {
	return &UploadDigitalKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("upload_digital_knowledge_source_doc", ctx)}
}

// UploadDigitalKnowledgeSourceDoc 把用户本地上传的 docx/xlsx 解析成纯正文，
// 以 source_type=file 在知识源的索引集合下建文档，由 platform 异步分段索引。
// 与在线文档（source_type=doc、引用云文档）不同，本地文件是「导入快照」：
// 正文已落在 platform，后续不会跟随某个在线源更新。
func (l *UploadDigitalKnowledgeSourceDocLogic) UploadDigitalKnowledgeSourceDoc(req *types.UploadDigitalKnowledgeSourceDocReq) (resp *types.UploadDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return nil, errors.New("知识源尚未就绪，请先保存知识源")
	}

	ext := strings.ToLower(path.Ext(req.FileName))
	if ext != ".docx" && ext != ".xlsx" {
		return &types.UploadDigitalKnowledgeSourceDocRes{
			Msg: "仅支持 .docx / .xlsx 文件",
		}, nil
	}

	data, decErr := base64.StdEncoding.DecodeString(req.Data)
	if decErr != nil {
		return &types.UploadDigitalKnowledgeSourceDocRes{Msg: "文件内容不是合法的 base64"}, nil
	}
	if len(data) == 0 {
		return &types.UploadDigitalKnowledgeSourceDocRes{Msg: "文件为空"}, nil
	}
	if len(data) > maxKnowledgeUploadBytes {
		return &types.UploadDigitalKnowledgeSourceDocRes{Msg: "文件超过 20MB 上限"}, nil
	}

	var text string
	if ext == ".docx" {
		text, err = parseKnowledgeDocx(data)
	} else {
		text, err = parseKnowledgeXlsx(data)
	}
	if err != nil {
		l.logger.Warn(model.LogMsg{Text: "本地知识文件解析失败", Data: map[string]interface{}{"file": req.FileName, "err": err.Error()}})
		return &types.UploadDigitalKnowledgeSourceDocRes{Msg: "文件解析失败：" + err.Error()}, nil
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return &types.UploadDigitalKnowledgeSourceDocRes{Msg: "未能从文件中解析出正文"}, nil
	}

	title := strings.TrimSpace(strings.TrimSuffix(path.Base(req.FileName), ext))
	if title == "" {
		title = "导入文档"
	}
	docType := strings.ToLower(strings.TrimSpace(req.DocType))
	if docType == "" {
		docType = "doc"
	}
	mediaType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if ext == ".xlsx" {
		mediaType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	created, createErr := l.svcCtx.PlatformRpc.CreateDocument(l.ctx, &platform.CreateDocumentReq{
		CollectionId: collectionID,
		Title:        title,
		MediaType:    mediaType,
		DocType:      docType,
		ContentText:  text,
		SourceType:   "file",
		SyncMode:     "import",
		RagEnabled:   true,
	})
	if createErr != nil {
		l.logger.Error(model.LogMsg{Text: "本地知识文件入库失败", Data: map[string]interface{}{"file": req.FileName, "err": createErr.Error()}})
		return nil, errors.New("文件入库失败")
	}
	return &types.UploadDigitalKnowledgeSourceDocRes{
		DocumentID: created.DocumentId,
		Title:      title,
		ByteSize:   int64(len([]byte(text))),
	}, nil
}
