package open

import (
	"context"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperListLogic {
	return &GetDeveloperListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDeveloperListLogic) GetDeveloperList(req *types.GetDeveloperListReq) (resp *types.GetDeveloperListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListDevelopers(l.ctx, &open_rpc.ListDevelopersReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Status:   int32(req.Status),
	})
	if err != nil {
		l.Errorf("获取开发者列表失败: %v", err)
		return nil, err
	}

	list := make([]types.DeveloperInfo, 0, len(rpcRes.List))
	for _, dev := range rpcRes.List {
		list = append(list, types.DeveloperInfo{
			ID:          fmt.Sprintf("%d", dev.Id),
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
			CreatedAt:   dev.CreatedAt,
		})
	}

	return &types.GetDeveloperListRes{Total: rpcRes.Total, List: list}, nil
}
