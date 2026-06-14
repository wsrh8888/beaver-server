package open

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AuditDeveloperLogic) AuditDeveloper(req *types.AuditDeveloperReq) (resp *types.AuditDeveloperRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.New("无效的申请ID")
	}

	_, err = l.svcCtx.OpenRpc.AuditDeveloper(l.ctx, &open_rpc.AuditDeveloperReq{
		Id:          id,
		Status:      int32(req.Status),
		AuditBy:     req.UserID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		l.Errorf("审核开发者失败: %v", err)
		return nil, err
	}

	return &types.AuditDeveloperRes{}, nil
}
