package developer

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperDetailLogic {
	return &GetDeveloperDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDeveloperDetailLogic) GetDeveloperDetail(req *types.GetDeveloperDetailReq) (resp *types.GetDeveloperDetailRes, err error) {
	_, ok := l.ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("未登录")
	}
	if req.ID == 0 {
		return nil, errors.New("id 不能为空")
	}

	rpcRes, err := l.svcCtx.OpenRpc.GetDeveloper(l.ctx, &open_rpc.GetDeveloperReq{Id: uint64(req.ID)})
	if err != nil {
		return nil, err
	}

	dev := rpcRes.Developer
	return &types.GetDeveloperDetailRes{
		Developer: types.DeveloperInfo{
			ID:          uint(dev.Id),
			UserID:      dev.UserId,
			RealName:    dev.RealName,
			CompanyName: dev.CompanyName,
			Phone:       dev.Phone,
			Email:       dev.Email,
			Description: dev.Description,
			Status:      int(dev.Status),
			AuditBy:     dev.AuditBy,
			AuditTime:   dev.AuditTime,
			AuditRemark: dev.AuditRemark,
			CreatedAt:   dev.CreatedAt / 1000,
		},
	}, nil
}
