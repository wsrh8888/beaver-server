package logic

import (
	"context"

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewLogic) Preview(req *types.PreviewReq) (resp *types.PreviewRes, err error) {
	// todo: add your logic here and delete this line

	return &types.PreviewRes{}, nil
}
