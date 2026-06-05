package open

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
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
	userID := req.ApplicantUserID
	if userID == "" {
		userID = req.UserID
	}
	if userID == "" {
		return nil, errors.New("未指定申请用户")
	}

	rpcRes, err := l.svcCtx.OpenRpc.ApplyDeveloper(l.ctx, &open_rpc.ApplyDeveloperReq{
		UserId:      userID,
		RealName:    req.RealName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("开发者申请失败: %v", err)
		return nil, err
	}

	return &types.ApplyDeveloperRes{ID: fmt.Sprintf("%d", rpcRes.Id)}, nil
}
