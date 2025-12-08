package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/device"
	"beaver/utils/jwts"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailPasswordLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailPasswordLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailPasswordLoginLogic {
	return &EmailPasswordLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailPasswordLoginLogic) EmailPasswordLogin(req *types.EmailPasswordLoginReq) (resp *types.EmailPasswordLoginRes, err error) {
	// 查询用户信息
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "email = ?", req.Email).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 调试日志：检查用户数据
	logx.Infof("登录用户信息: UserID=%s, NickName=%s, Email=%s", user.UserID, user.NickName, user.Email)

	// 验证密码
	if !pwd.CheckPad(user.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 根据User-Agent识别设备类型
	userAgent := l.ctx.Value("user-agent")
	fmt.Println("user-agent:", userAgent)

	var deviceType string
	if userAgent == nil {
		deviceType = "unknown" // 如果没有UA信息，使用默认类型
		logx.Info("无法识别设备类型，使用默认类型：unknown")
	} else {
		deviceType = device.GetDeviceType(userAgent.(string))
	}

	// 生成token，包含设备信息
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		UserID:   user.UserID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 检查是否已有其他设备登录
	key := fmt.Sprintf("login_%s_%s", user.UserID, deviceType)
	oldLoginInfo, err := l.svcCtx.Redis.Get(key).Result()
	if err == nil && oldLoginInfo != "" {
		// 解析旧登录信息
		var loginInfo map[string]interface{}
		if err := json.Unmarshal([]byte(oldLoginInfo), &loginInfo); err == nil {
			oldDeviceID := loginInfo["device_id"].(string)
			if oldDeviceID != req.DeviceID {
				fmt.Println("不是同一设备，需要通知踢出旧设备:", oldDeviceID)
				// 不是同一设备，需要通知踢出旧设备
				// notifyForceOffline(user.UserID, oldDeviceID)
			}
		}
	}

	// 存储新的登录信息（包含token和设备信息）
	loginInfo := map[string]interface{}{
		"token":       token,
		"device_id":   req.DeviceID,
		"device_type": deviceType,
		"login_time":  time.Now().Format("2006-01-02 15:04:05"),
		"user_agent":  userAgent,
	}

	loginInfoJson, _ := json.Marshal(loginInfo)
	err = l.svcCtx.Redis.Set(key, string(loginInfoJson), time.Hour*48).Err()
	if err != nil {
		logx.Errorf("存储登录信息失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	return &types.EmailPasswordLoginRes{
		Token:  token,
		UserID: user.UserID,
	}, nil
}
