package developer

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(req *types.ApplyDeveloperReq) (resp *types.ApplyDeveloperRes, err error) {
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录或登录已过期")
	}

	_, err = l.svcCtx.OpenRpc.ApplyDeveloper(l.ctx, &open_rpc.ApplyDeveloperReq{
		UserId:      userID,
		RealName:    req.RealName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Description: req.Description,
	})
	if err != nil {
		logx.Errorf("开发者申请失败: %v", err)
		return nil, err
	}

	return &types.ApplyDeveloperRes{}, nil
}
