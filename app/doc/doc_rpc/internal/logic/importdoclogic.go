package logic

import (
	"context"
	"fmt"
	"path"
	"strings"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ImportDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewImportDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportDocLogic {
	return &ImportDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ImportDocLogic", ctx),
	}
}

// ImportDoc 导入 docx/pdf 为文档域权威文档：解析正文 -> 建 page 节点 + 文档 -> 写正文。
// pdf 依赖（pdftotext）缺失时优雅降级：仍建文档、仅登记不解析正文（parsed=false）。
func (l *ImportDocLogic) ImportDoc(in *doc_rpc.ImportDocReq) (*doc_rpc.ImportDocRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}
	if len(in.Data) == 0 {
		return nil, errInvalidArg
	}
	maxBytes := l.svcCtx.Config.Import.MaxBytes
	if maxBytes <= 0 {
		maxBytes = 20971520
	}
	if int64(len(in.Data)) > maxBytes {
		return nil, fmt.Errorf("doc: 文件超过导入上限 %d 字节", maxBytes)
	}

	// 按后缀判定文件类型并解析正文。
	ext := strings.ToLower(path.Ext(in.FileName))
	fileKind := "unsupported"
	var text string
	var parsed bool
	var msg string
	switch ext {
	case ".docx":
		fileKind = "docx"
		t, perr := parseDocx(in.Data)
		if perr != nil {
			msg = "docx 解析失败：" + perr.Error()
		} else {
			text = t
			parsed = true
		}
	case ".pdf":
		fileKind = "pdf"
		t, ok, m := parsePdf(l.svcCtx.Config.Import.PdfBinPath, in.Data)
		text = t
		parsed = ok
		msg = m
	default:
		msg = "不支持的文件类型，仅登记未解析正文"
	}

	title := trim(in.Title)
	if title == "" {
		base := path.Base(in.FileName)
		title = trim(strings.TrimSuffix(base, path.Ext(base)))
	}
	if title == "" {
		title = "导入文档"
	}

	// 建 page 节点 + 文档（骨架正文先落库，保证节点/文档成对）。
	nodeID := genID("node")
	docID := genID("doc")
	doc := doc_models.DocDocument{
		DocID:       docID,
		NodeID:      nodeID,
		SpaceID:     space.SpaceID,
		Title:       title,
		ContentJSON: skeletonContentJSON,
		Version:     1,
		OwnerUserID: trim(in.OperatorUserId),
		Status:      1,
	}
	doc.ContentDigest = sha256hex(doc.ContentText)
	node := doc_models.DocNode{
		NodeID:   nodeID,
		SpaceID:  space.SpaceID,
		ParentID: trim(in.ParentNodeId),
		Title:    title,
		Type:     "page",
		DocID:    docID,
	}
	if err := l.svcCtx.DB.Create(&doc).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "导入创建文档失败", Data: err})
		return nil, err
	}
	if err := l.svcCtx.DB.Create(&node).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "导入创建节点失败", Data: err})
		return nil, err
	}

	// 解析成功则写入正文 + 生成版本（复用 applyDocContent，失败不阻断导入结果返回）。
	if parsed && text != "" {
		contentJSON := textToContentJSON(text)
		if _, aerr := applyDocContent(l.ctx, l.svcCtx, &doc, space, contentJSON, text, trim(in.OperatorUserId)); aerr != nil {
			l.logger.Error(model.LogMsg{Text: "导入写入正文失败", Data: map[string]any{"docId": docID, "err": aerr.Error()}})
		}
	}

	return &doc_rpc.ImportDocRes{
		DocId:       docID,
		NodeId:      nodeID,
		ContentText: text,
		FileKind:    fileKind,
		Parsed:      parsed,
		Message:     msg,
	}, nil
}
