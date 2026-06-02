package developer

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperDetailLogic {
	return &GetDeveloperDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeveloperDetailLogic) GetDeveloperDetail(req *types.GetDeveloperDetailReq) (resp *types.GetDeveloperDetailRes, err error) {
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := l.svcCtx.RequireDeveloper(userID); err != nil {
		return nil, err
	}

	// TODO: 实现开发者详情查询逻辑
	return nil, nil
}
