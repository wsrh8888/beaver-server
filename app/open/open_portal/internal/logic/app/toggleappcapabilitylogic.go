package app

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

type ToggleAppCapabilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用/禁用应用能力（对标飞书）
func NewToggleAppCapabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleAppCapabilityLogic {
	return &ToggleAppCapabilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleAppCapabilityLogic) ToggleAppCapability(req *types.ToggleAppCapabilityReq) (resp *types.ToggleAppCapabilityRes, err error) {

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 根据能力类型更新对应的开关
	var enabled bool
	switch req.Capability {
	case "robot":
		if req.Enable {
			app.EnableRobot = 1
			app.EnableWebhook = 1
			enabled = true
		} else {
			app.EnableRobot = 0
			enabled = false
		}
	case "oauth":
		if req.Enable {
			app.EnableOAuth = 1
			enabled = true
		} else {
			app.EnableOAuth = 0
			enabled = false
		}
	case "webhook":
		if req.Enable {
			app.EnableWebhook = 1
			enabled = true
		} else {
			app.EnableWebhook = 0
			enabled = false
		}
	default:
		return nil, errors.New("不支持的能力类型")
	}

	// 3. 保存更新
	if err := l.svcCtx.DB.Save(&app).Error; err != nil {
		logx.Errorf("更新应用能力失败: %v", err)
		return nil, errors.New("更新应用能力失败")
	}

	if req.Capability == "robot" && req.Enable {
		if err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
			logx.Errorf("创建 Robot 用户失败: app_id=%s err=%v", req.AppID, err)
			return nil, errors.New("启用 Robot 成功，但创建 IM 用户失败，请稍后重试")
		}
	}

	logx.Infof("应用 %s 的 %s 能力已%s", req.AppID, req.Capability, map[bool]string{true: "启用", false: "禁用"}[req.Enable])

	return &types.ToggleAppCapabilityRes{
		Enabled: enabled,
	}, nil
}

func ensurePortalAppRobot(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) error {
	var robot open_models.OpenAppRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
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
		return fmt.Errorf("user create: %w", err)
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
	}
	return db.Save(&robot).Error
}
