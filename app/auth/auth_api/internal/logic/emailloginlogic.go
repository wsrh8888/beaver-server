package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/middleware/ua"
	"beaver/utils/device"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailLoginLogic) EmailLogin(req *types.EmailLoginReq) (resp *types.EmailLoginRes, err error) {
	// 验证验证码
	err = l.verifyEmailCode(req.Email, req.Code, "login")
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "email = ?", req.Email).Error
	if err != nil {
		// 如果用户不存在，自动注册
		user, err = l.createUserByEmail(req.Email)
		if err != nil {
			return nil, errors.New("用户注册失败")
		}
	}

	// 1. 获取精准设备信息（由 UA 中间件预处理注入）
	preciseType, _ := l.ctx.Value(ua.KeyDeviceType).(string)
	deviceGroup, _ := l.ctx.Value(ua.KeyDeviceGroup).(string)

	if preciseType == "" || preciseType == device.DeviceUnknown {
		logx.WithContext(l.ctx).Errorf("非法请求：无法解析设备类型")
		return nil, errors.New("不支持的设备类型")
	}

	// 核心安全校验：拒绝空指纹设备登录
	if req.DeviceID == "" {
		logx.WithContext(l.ctx).Errorf("连接拒绝：缺少物理指纹, 用户: %s", req.Email)
		return nil, errors.New("无法识别的物理设备，请联系管理员")
	}

	key := fmt.Sprintf("user_authentication_session:%s:%s", user.UserID, deviceGroup)

	// 3. 生成 token，包含基本信息 + 设备唯一标识 (DeviceID/GUID)
	// 实现“设备绑定”：即使 Token 泄露，非指纹设备也无法使用
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		UserID:   user.UserID,
		DeviceID: req.DeviceID, // 将 GUID 写入 Token 荷载
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 4. 检查同槽位是否有其他设备在线（互踢逻辑）
	oldLoginInfo, err := l.svcCtx.Redis.Get(key).Result()
	if err == nil && oldLoginInfo != "" {
		var oldInfo map[string]interface{}
		if err := json.Unmarshal([]byte(oldLoginInfo), &oldInfo); err == nil {
			oldDeviceID, _ := oldInfo["device_id"].(string)
			oldPreciseType, _ := oldInfo["device_type"].(string)

			// 如果 DeviceID 不同，说明是同族的不同物理设备登录，需要执行“顶号”
			if oldDeviceID != req.DeviceID {
				logx.WithContext(l.ctx).Infof("【系统通知】用户 %s 在新设备 %s (%s) 登录，正在踢出旧设备 %s (%s)",
					user.UserID, req.DeviceID, preciseType, oldDeviceID, oldPreciseType)
				// TODO: 这里后续可以接入 WS 推送，实时通知旧设备“你已被踢下线”
			}
		}
	}

	// 5. 存储新的登录信息（槽位占位）
	loginInfo := map[string]interface{}{
		"token":        token,
		"device_id":    req.DeviceID,
		"device_type":  preciseType, // 存入精准类型，方便展示和定位
		"device_group": deviceGroup,
		"login_time":   time.Now().Format("2006-01-02 15:04:05"),
	}

	loginInfoJson, _ := json.Marshal(loginInfo)
	err = l.svcCtx.Redis.Set(key, string(loginInfoJson), time.Hour*time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Err()
	if err != nil {
		logx.Errorf("存储登录信息失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	return &types.EmailLoginRes{
		Token:  token,
		UserID: user.UserID,
	}, nil
}

// 通过邮箱创建新用户
func (l *EmailLoginLogic) createUserByEmail(email string) (user_models.UserModel, error) {
	// 生成随机昵称
	nickName := fmt.Sprintf("用户%s", email[:strings.Index(email, "@")])

	// 创建用户
	user := user_models.UserModel{
		Email:    email,
		NickName: nickName,
		UserID:   generateUUID(),
		// 其他字段使用默认值
	}

	err := l.svcCtx.DB.Create(&user).Error
	if err != nil {
		return user_models.UserModel{}, err
	}

	return user, nil
}

// 生成UUID
func generateUUID() string {
	// 这里可以使用现有的UUID生成工具
	// 暂时使用简单的实现
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

// 验证邮箱验证码
func (l *EmailLoginLogic) verifyEmailCode(email, code, codeType string) error {
	// 从Redis获取存储的验证码
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}

	// 验证验证码
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码
	l.svcCtx.Redis.Del(codeKey)

	return nil
}
