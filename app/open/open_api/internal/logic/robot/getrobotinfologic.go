package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRobotInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRobotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotInfoLogic {
	return &GetRobotInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetRobotInfoLogic) GetRobotInfo(authorization string) (resp *types.GetRobotInfoRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppCapability(app, true, false); err != nil {
		return nil, err
	}

	robot, err := utils.EnsureAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	return &types.GetRobotInfoRes{
		RobotID:   robot.RobotID,
		RobotName: robot.RobotName,
		Avatar:    robot.Avatar,
		AppID:     app.AppID,
	}, nil
}
