package developer

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
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
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录")
	}
	if req.Status != 1 && req.Status != 2 {
		return nil, fmt.Errorf("无效的状态值")
	}

	_, err = l.svcCtx.OpenRpc.AuditDeveloper(l.ctx, &open_rpc.AuditDeveloperReq{
		Id:          uint64(req.ID),
		Status:      int32(req.Status),
		AuditBy:     userID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		return nil, err
	}

	l.Infof("开发者申请 %d 已审核，状态: %d", req.ID, req.Status)
	return &types.AuditDeveloperRes{}, nil
}
