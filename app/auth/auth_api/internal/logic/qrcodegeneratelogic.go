package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

// 二维码状态常量（供三个 logic 共用，同一 package 可直接访问）
const (
	QrcodeStatusPending   = "pending"   // 等待移动端扫码
	QrcodeStatusConfirmed = "confirmed" // 移动端已确认，PC 端可取 JWT
	QrcodeStatusExpired   = "expired"   // 已过期或已被使用（一次性）

	qrcodeTTL    = 3 * time.Minute // 二维码本身的有效期
	qrcodeKeyFmt = "qrcode:%s"     // Redis key 格式
)

// qrcodeTokenExpireHours 扫码登录固定 12 小时，source 仅作来源记录不影响时长
const qrcodeTokenExpireHours = 12

// QrcodeSession 存入 Redis 的二维码会话，供 scan/status 两个 logic 读写
type QrcodeSession struct {
	Status        string `json:"status"`
	Source        string `json:"source"`                  // generate 时传入的来源标识，透传至最终响应
	ScannedUserID string `json:"scannedUserId,omitempty"` // scan 后写入的扫码人 userID
}

type QrcodeGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 扫码登录-生成二维码（PC 端调用）
func NewQrcodeGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeGenerateLogic {
	return &QrcodeGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeGenerateLogic) QrcodeGenerate(req *types.QrcodeGenerateReq) (resp *types.QrcodeGenerateRes, err error) {
	source := req.Source
	if source == "" {
		source = "login"
	}

	// 生成 UUID v4 作为二维码唯一 token
	token := util.NewV4().String()

	session := QrcodeSession{
		Status: QrcodeStatusPending,
		Source: source,
	}
	sessionJSON, _ := json.Marshal(session)
	key := fmt.Sprintf(qrcodeKeyFmt, token)

	if err = l.svcCtx.Redis.Set(key, string(sessionJSON), qrcodeTTL).Err(); err != nil {
		logx.Errorf("qrcode generate: redis set failed key=%s err=%v", key, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	return &types.QrcodeGenerateRes{
		Token:    token,
		ExpireAt: time.Now().Add(qrcodeTTL).Unix(),
	}, nil
}
