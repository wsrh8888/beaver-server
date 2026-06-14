package robot

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetRobotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRobotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotConfigLogic {
	return &GetRobotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRobotConfigLogic) GetRobotConfig(req *types.GetRobotConfigReq) (resp *types.GetRobotConfigRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}
	if app.EnableRobot != 1 {
		return nil, errors.New("应用未启用智能机器人能力，请先在应用能力中开启 robot")
	}

	robot, err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app)
	if err != nil {
		return nil, errors.New("获取 Robot 配置失败")
	}

	return &types.GetRobotConfigRes{
		Config: robotConfigFromModel(req.AppID, robot),
	}, nil
}

func robotConfigFromModel(appID string, robot *open_models.OpenAppRobot) types.RobotConfigInfo {
	return types.RobotConfigInfo{
		AppID:            appID,
		RobotID:          robot.RobotID,
		RobotName:        robot.RobotName,
		Avatar:           robot.Avatar,
		WelcomeMessage:   robot.WelcomeMessage,
		CommandPrefix:    robot.CommandPrefix,
		EnableSingleChat: robot.EnableSingleChat == 1,
		EnableGroupChat:  robot.EnableGroupChat == 1,
		EnableAtMention:  robot.EnableAtMention == 1,
		Status:           robot.Status,
	}
}

func ensurePortalAppRobot(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) (*open_models.OpenAppRobot, error) {
	var robot open_models.OpenAppRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		return &robot, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nickName := app.Name
	if nickName == "" {
		nickName = "Robot"
	}

	createRes, err := userRpc.UserCreate(ctx, &user.UserCreateReq{
		NickName: nickName,
		UserType: int32(user_models.UserTypeRobot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}

	robot = open_models.OpenAppRobot{
		AppID:            app.AppID,
		RobotID:          createRes.UserID,
		RobotName:        nickName,
		Avatar:           app.Icon,
		Status:           1,
		EnableSingleChat: 1,
		EnableGroupChat:  1,
		EnableAtMention:  1,
		CommandPrefix:    "/",
	}
	if err := db.Save(&robot).Error; err != nil {
		return nil, err
	}
	return &robot, nil
}