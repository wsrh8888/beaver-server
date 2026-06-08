package moderation

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type ExecuteUserControlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewExecuteUserControlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteUserControlLogic {
	return &ExecuteUserControlLogic{logger: logger.New("execute_user_control"), ctx: ctx, svcCtx: svcCtx}
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
		logx.WithContext(l.ctx).Errorf("用户管控失败: %v", err)
		return nil, err
	}
	l.logger.Info(model.LogMsg{
		Text: "用户管控执行成功",
		Data: map[string]interface{}{
			"operatorId": req.OperatorID,
			"userId":     req.UserID,
			"action":     req.Action,
		},
	})
	return &types.ExecuteUserControlRes{}, nil
}
