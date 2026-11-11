package logic

import (
	"context"
	"encoding/base64"
	"strings"

	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ExportDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewExportDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportDocLogic {
	return &ExportDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ExportDocLogic", ctx),
	}
}

func (l *ExportDocLogic) ExportDoc(in *doc_rpc.ExportDocReq) (*doc_rpc.ExportDocRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.OperatorUserId, roleViewer)
	if err != nil {
		return nil, err
	}
	format := strings.ToLower(trim(in.Format))
	if format == "" {
		format = "md"
	}
	name := strings.TrimSpace(doc.Title)
	if name == "" {
		name = "untitled"
	}

	var data, mime, fileName string
	switch format {
	case "md":
		data, mime, fileName = pmToMarkdown(doc.ContentJSON), "text/markdown", name+".md"
	case "html":
		data, mime, fileName = pmToHTML(doc.ContentJSON), "text/html", name+".html"
	case "txt":
		data, mime, fileName = doc.ContentText, "text/plain", name+".txt"
	case "json":
		data, mime, fileName = doc.ContentJSON, "application/json", name+".json"
	case "docx":
		b, derr := buildDocx(doc.Title, doc.ContentText)
		if derr != nil {
			return nil, derr
		}
		data = base64.StdEncoding.EncodeToString(b)
		mime = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		fileName = name + ".docx"
	default:
		return nil, errInvalidArg
	}
	res := &doc_rpc.ExportDocRes{FileName: fileName, MimeType: mime, Data: data}
	if format != "docx" {
		res.Text = data
	}
	auditJSON(l.svcCtx.DB, in.OperatorUserId, "doc.export", "doc", doc.DocID,
		map[string]any{"format": format})
	return res, nil
}
