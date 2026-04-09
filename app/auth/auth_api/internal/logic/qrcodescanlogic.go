package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/jwts"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type QrcodeScanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 扫码登录-确认扫码（移动端调用，需登录）
func NewQrcodeScanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeScanLogic {
	return &QrcodeScanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeScanLogic) QrcodeScan(req *types.QrcodeScanReq) (resp *types.QrcodeScanRes, err error) {
	// 1. 解析移动端 JWT，获取扫码人身份
	claims, err := jwts.ParseToken(req.AuthToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("身份验证失败，请重新登录")
	}

	// 2. 查询二维码会话
	key := fmt.Sprintf(qrcodeKeyFmt, req.Token)
	sessionStr, err := l.svcCtx.Redis.Get(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("二维码已过期，请刷新后重试")
	}
	if err != nil {
		logx.Errorf("qrcode scan: redis get failed key=%s err=%v", key, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	// 3. 检查当前状态
	var session QrcodeSession
	if err = json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}
	if session.Status != QrcodeStatusPending {
		return nil, fmt.Errorf("二维码已被使用或已过期")
	}

	// 4. 更新状态为 confirmed，写入扫码人 userID，保持原 TTL 不变
	ttl := l.svcCtx.Redis.TTL(key).Val()
	session.Status = QrcodeStatusConfirmed
	session.ScannedUserID = claims.UserID
	updatedJSON, _ := json.Marshal(session)

	if err = l.svcCtx.Redis.Set(key, string(updatedJSON), ttl).Err(); err != nil {
		logx.Errorf("qrcode scan: redis update failed key=%s err=%v", key, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	return &types.QrcodeScanRes{}, nil
}
