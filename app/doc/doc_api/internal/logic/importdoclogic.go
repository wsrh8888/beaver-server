package logic

import (
	"context"
	"encoding/base64"
	"errors"

	"beaver/app/doc/doc_api/internal/svc"
	"beaver/app/doc/doc_api/internal/types"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"

	"github.com/zeromicro/go-zero/core/logx"

	beaverlog "beaver/utils/beaverlog"
)

// goctl 会把 logx 写进 imports，占位避免未使用导入。
var _ = logx.WithContext

// errInvalidBase64 导入文件 data 字段非合法 base64。
var errInvalidBase64 = errors.New("doc: 导入数据非合法 base64")

type ImportDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 导入 docx/pdf 文件为文档
func NewImportDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportDocLogic {
	return &ImportDocLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ImportDoc", ctx),
	}
}

func (l *ImportDocLogic) ImportDoc(req *types.ImportDocReq) (resp *types.ImportDocRes, err error) {
	data, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, errInvalidBase64
	}
	res, err := l.svcCtx.Doc.ImportDoc(l.ctx, &doc_rpc.ImportDocReq{
		SpaceId:        req.SpaceId,
		ParentNodeId:   req.ParentNodeId,
		Title:          req.Title,
		FileName:       req.FileName,
		Data:           data,
		OperatorUserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ImportDocRes{
		DocId:       res.DocId,
		NodeId:      res.NodeId,
		ContentText: res.ContentText,
		FileKind:    res.FileKind,
		Parsed:      res.Parsed,
		Message:     res.Message,
	}, nil
}
