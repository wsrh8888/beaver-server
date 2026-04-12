package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/device"
	"beaver/utils/jwts"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type QrcodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 扫码登录-轮询状态（PC 端调用）
func NewQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeStatusLogic {
	return &QrcodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeStatusLogic) QrcodeStatus(req *types.QrcodeStatusReq) (resp *types.QrcodeStatusRes, err error) {
	// 1. 查询二维码会话
	key := fmt.Sprintf(qrcodeKeyFmt, req.Token)
	sessionStr, err := l.svcCtx.Redis.Get(key).Result()
	if err == redis.Nil {
		return &types.QrcodeStatusRes{Status: QrcodeStatusExpired}, nil
	}
	if err != nil {
		logx.Errorf("qrcode status: redis get failed key=%s err=%v", key, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	var session QrcodeSession
	if err = json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}

	// 2. 未扫码，继续等待
	if session.Status == QrcodeStatusPending {
		return &types.QrcodeStatusRes{Status: QrcodeStatusPending}, nil
	}

	// 3. 已确认：生成 PC 端 JWT，写入登录态，删除一次性 session
	if session.Status == QrcodeStatusConfirmed {
		var user user_models.UserModel
		if err = l.svcCtx.DB.Take(&user, "user_id = ?", session.ScannedUserID).Error; err != nil {
			logx.Errorf("qrcode status: user not found userID=%s err=%v", session.ScannedUserID, err)
			return nil, fmt.Errorf("用户不存在")
		}

		// JWT 有效期固定 12h，source 仅作来源记录
		jwtExpireHours := qrcodeTokenExpireHours

		token, err := jwts.GenToken(jwts.JwtPayLoad{
			NickName: user.NickName,
			UserID:   user.UserID,
		}, l.svcCtx.Config.Auth.AccessSecret, jwtExpireHours)
		if err != nil {
			logx.Errorf("qrcode status: gen token failed userID=%s err=%v", user.UserID, err)
			return nil, fmt.Errorf("服务内部异常")
		}

		// 识别 PC 设备类型（与普通登录逻辑完全一致）
		deviceType := "desktop"
		if ua := l.ctx.Value("user-agent"); ua != nil {
			deviceType = device.GetDeviceType(ua.(string))
		}

		// 写入登录态（与 emailPasswordLoginLogic 保持完全一致的 key 格式和数据结构）
		loginKey := fmt.Sprintf("login_%s_%s", user.UserID, deviceType)
		loginInfo := map[string]any{
			"token":       token,
			"device_id":   req.DeviceID,
			"device_type": deviceType,
			"login_time":  time.Now().Format("2006-01-02 15:04:05"),
			"source":      session.Source,
		}
		loginInfoJSON, _ := json.Marshal(loginInfo)
		loginTTL := time.Duration(jwtExpireHours) * time.Hour
		if err = l.svcCtx.Redis.Set(loginKey, string(loginInfoJSON), loginTTL).Err(); err != nil {
			logx.Errorf("qrcode status: set login info failed key=%s err=%v", loginKey, err)
			return nil, fmt.Errorf("服务内部异常")
		}

		// 删除二维码 session，防止重复使用
		l.svcCtx.Redis.Del(key)

		return &types.QrcodeStatusRes{
			Status: QrcodeStatusConfirmed,
			Token:  token,
			UserID: user.UserID,
			Source: session.Source,
		}, nil
	}

	return &types.QrcodeStatusRes{Status: QrcodeStatusExpired}, nil
}
