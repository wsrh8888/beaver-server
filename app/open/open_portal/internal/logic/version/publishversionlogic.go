package version

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布版本
func NewPublishVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVersionLogic {
	return &PublishVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishVersionLogic) PublishVersion(req *types.PublishVersionReq) (resp *types.PublishVersionRes, err error) {

	return nil, errors.New("功能暂未开放")
}
