package moderation

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteUserControlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteUserControlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteUserControlLogic {
	return &ExecuteUserControlLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ExecuteUserControlLogic) ExecuteUserControl(req *types.ExecuteUserControlReq) (resp *types.ExecuteUserControlRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.Action != "ban_user" && req.Action != "unban_user" {
		return nil, errors.New("仅支持 ban_user / unban_user")
	}

	act := types.ModerationControlAction{
		Action: req.Action,
		Target: req.UserID,
		Reason: req.Reason,
	}
	if err = executeControlAction(l.ctx, l.svcCtx, req.OperatorID, req.CaseID, act); err != nil {
		l.Errorf("用户管控失败: %v", err)
		return nil, err
	}
	return &types.ExecuteUserControlRes{}, nil
}
