package version

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitVersionReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交版本审核
func NewSubmitVersionReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitVersionReviewLogic {
	return &SubmitVersionReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitVersionReviewLogic) SubmitVersionReview(req *types.SubmitVersionReviewReq) (resp *types.SubmitVersionReviewRes, err error) {
	return nil, errors.New("功能暂未开放")
}
