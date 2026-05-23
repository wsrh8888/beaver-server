package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
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

	// 查询用户凭证并验证密码
	var credential auth_models.AuthCredentialModel
	err = l.svcCtx.DB.Take(&credential, "user_id = ?", user.UserID).Error
	if err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	if !pwd.CheckPad(credential.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 1. 获取 User-Agent
	userAgent, _ := l.ctx.Value("user-agent").(string)
	preciseType := device.GetDeviceType(userAgent)

	if preciseType == "" || preciseType == device.DeviceUnknown {
		logx.WithContext(l.ctx).Errorf("非法请求：无法解析设备类型, UA: %s", userAgent)
		return nil, errors.New("不支持的设备类型")
	}

	// 2. 核心安全校验：拒绝空指纹设备登录
	if req.DeviceID == "" {
		logx.WithContext(l.ctx).Errorf("连接拒绝：缺少物理指纹, 用户: %s", req.Email)
		return nil, errors.New("无法识别的物理设备，请联系管理员")
	}

	// 3. 确定登录槽位 (Group: desktop/mobile)
	deviceGroup := device.GetDeviceGroup(preciseType)
	key := fmt.Sprintf("user_authentication_session:%s:%s", user.UserID, deviceGroup)

	// 4. 生成 token，包含基本信息 + 设备唯一标识 (DeviceID/GUID)
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		UserID:   user.UserID,
		DeviceID: req.DeviceID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 5. 检查同槽位是否有其他设备在线（互踢逻辑）
	oldLoginInfo, err := l.svcCtx.Redis.Get(key).Result()
	if err == nil && oldLoginInfo != "" {
		var oldInfo map[string]interface{}
		if err := json.Unmarshal([]byte(oldLoginInfo), &oldInfo); err == nil {
			oldDeviceID, _ := oldInfo["device_id"].(string)
			if oldDeviceID != req.DeviceID {
				logx.WithContext(l.ctx).Infof("【系统通知】用户 %s 在新设备 %s (%s) 登录，正在踢出旧设备 %s",
					user.UserID, req.DeviceID, preciseType, oldDeviceID)
			}
		}
	}

	// 6. 存储新的登录信息（槽位占位）
	loginInfo := map[string]interface{}{
		"token":        token,
		"device_id":    req.DeviceID,
		"device_type":  preciseType,
		"device_group": deviceGroup,
		"login_time":   time.Now().Format("2006-01-02 15:04:05"),
		"user_agent":   userAgent,
	}

	loginInfoJson, _ := json.Marshal(loginInfo)
	err = l.svcCtx.Redis.Set(key, string(loginInfoJson), time.Hour*time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Err()
	if err != nil {
		logx.Errorf("存储登录信息失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 7. 更新登录记录
	now := time.Now()
	credential.LastLoginAt = &now
	credential.LoginCount++
	l.svcCtx.DB.Save(&credential)

	return &types.EmailPasswordLoginRes{
		Token:  token,
		UserID: user.UserID,
	}, nil
}
