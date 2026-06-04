package logic

import (
	"context"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRobotByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRobotByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotByUserIDLogic {
	return &GetRobotByUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRobotByUserIDLogic) GetRobotByUserID(in *open_rpc.GetRobotByUserIDReq) (*open_rpc.GetRobotByUserIDRes, error) {
	if in.RobotUserId == "" {
		return &open_rpc.GetRobotByUserIDRes{Found: false}, nil
	}

	var robot open_models.OpenAppRobot
	if err := l.svcCtx.DB.Where("robot_user_id = ? AND status = 1", in.RobotUserId).First(&robot).Error; err != nil {
		return &open_rpc.GetRobotByUserIDRes{Found: false}, nil
	}

	return &open_rpc.GetRobotByUserIDRes{
		Found:            true,
		AppId:            robot.AppID,
		RobotUserId:      robot.RobotID,
		EnableSingleChat: robot.EnableSingleChat == 1,
		EnableGroupChat:  robot.EnableGroupChat == 1,
		EnableAtMention:  robot.EnableAtMention == 1,
	}, nil
}
