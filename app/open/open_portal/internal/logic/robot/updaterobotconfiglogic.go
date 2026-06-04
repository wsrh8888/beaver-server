package robot

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRobotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRobotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRobotConfigLogic {
	return &UpdateRobotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRobotConfigLogic) UpdateRobotConfig(req *types.UpdateRobotConfigReq) (resp *types.UpdateRobotConfigRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}
	if app.EnableRobot != 1 {
		return nil, errors.New("应用未启用智能机器人能力")
	}

	robot, err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app)
	if err != nil {
		return nil, errors.New("更新 Robot 配置失败")
	}

	updates := map[string]interface{}{}
	if req.RobotName != "" {
		updates["robot_name"] = req.RobotName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.WelcomeMessage != "" {
		updates["welcome_message"] = req.WelcomeMessage
	}
	if req.CommandPrefix != "" {
		updates["command_prefix"] = req.CommandPrefix
	}
	if req.EnableSingleChat != nil {
		updates["enable_single_chat"] = boolToTinyInt(*req.EnableSingleChat)
	}
	if req.EnableGroupChat != nil {
		updates["enable_group_chat"] = boolToTinyInt(*req.EnableGroupChat)
	}
	if req.EnableAtMention != nil {
		updates["enable_at_mention"] = boolToTinyInt(*req.EnableAtMention)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(robot).Updates(updates).Error; err != nil {
			return nil, errors.New("更新 Robot 配置失败")
		}
	}

	if req.RobotName != "" || req.Avatar != "" {
		displayReq := &user.UserUpdateDisplayReq{UserId: robot.RobotID}
		if req.RobotName != "" {
			displayReq.NickName = req.RobotName
		}
		if req.Avatar != "" {
			displayReq.Avatar = req.Avatar
		}
		if _, err := l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, displayReq); err != nil {
			l.Errorf("同步 Robot IM 展示信息失败: robot=%s err=%v", robot.RobotID, err)
			return nil, errors.New("Robot 配置已保存，但同步 IM 昵称/头像失败")
		}
	}

	return &types.UpdateRobotConfigRes{}, nil
}

func boolToTinyInt(v bool) int {
	if v {
		return 1
	}
	return 0
}