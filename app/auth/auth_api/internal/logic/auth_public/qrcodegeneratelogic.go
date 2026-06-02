package auth_public

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

const (
	QrcodeStatusPending   = "pending"
	QrcodeStatusConfirmed = "confirmed"
	QrcodeStatusExpired   = "expired"

	QrcodeTTL            = 3 * time.Minute
	QrcodeKeyFmt         = "qrcode:%s"
	QrcodeTokenExpireHours = 12
)

type QrcodeSession struct {
	Status        string `json:"status"`
	Source        string `json:"source"`
	ScannedUserID string `json:"scannedUserId,omitempty"`
}

type QrcodeGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQrcodeGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeGenerateLogic {
	return &QrcodeGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeGenerateLogic) QrcodeGenerate(req *types.QrcodeGenerateReq) (*types.QrcodeGenerateRes, error) {
	source := req.Source
	if source == "" {
		source = "login"
	}

	token := util.NewV4().String()
	session := QrcodeSession{
		Status: QrcodeStatusPending,
		Source: source,
	}
	sessionJSON, _ := json.Marshal(session)
	key := fmt.Sprintf(QrcodeKeyFmt, token)

	if err := l.svcCtx.Redis.Set(key, string(sessionJSON), QrcodeTTL).Err(); err != nil {
		logx.Errorf("qrcode generate: redis set failed key=%s err=%v", key, err)
		return nil, fmt.Errorf("服务内部异常")
	}

	return &types.QrcodeGenerateRes{
		Token:    token,
		ExpireAt: time.Now().Add(QrcodeTTL).Unix(),
	}, nil
}
