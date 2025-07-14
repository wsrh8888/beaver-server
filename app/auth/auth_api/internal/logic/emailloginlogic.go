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

	// 根据User-Agent识别设备类型
	userAgent := l.ctx.Value("user-agent")
	fmt.Println("user-agent:", userAgent)

	var deviceType string
	if userAgent == nil {
		deviceType = "unknown"
		logx.Info("无法识别设备类型，使用默认类型：unknown")
	} else {
		deviceType = device.GetDeviceType(userAgent.(string))
	}

	// 生成token，包含设备信息
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Nickname: user.NickName,
		UserID:   user.UUID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 检查是否已有其他设备登录
	key := fmt.Sprintf("login_%s_%s", user.UUID, deviceType)
	oldLoginInfo, err := l.svcCtx.Redis.Get(key).Result()
	if err == nil && oldLoginInfo != "" {
		// 解析旧登录信息
		var loginInfo map[string]interface{}
		if err := json.Unmarshal([]byte(oldLoginInfo), &loginInfo); err == nil {
			oldDeviceID := loginInfo["device_id"].(string)
			if oldDeviceID != req.DeviceID {
				fmt.Println("不是同一设备，需要通知踢出旧设备:", oldDeviceID)
				// 不是同一设备，需要通知踢出旧设备
				// notifyForceOffline(user.UUID, oldDeviceID)
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

	return &types.EmailLoginRes{
		Token:  token,
		UserID: user.UUID,
	}, nil
}

// 通过邮箱创建新用户
func (l *EmailLoginLogic) createUserByEmail(email string) (user_models.UserModel, error) {
	// 生成随机昵称
	nickname := fmt.Sprintf("用户%s", email[:strings.Index(email, "@")])

	// 创建用户
	user := user_models.UserModel{
		Email:    email,
		NickName: nickname,
		UUID:     generateUUID(),
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
