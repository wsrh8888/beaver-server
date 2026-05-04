package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateQrCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 生成扫码登录二维码
func NewGenerateQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrCodeLogic {
	return &GenerateQrCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateQrCodeLogic) GenerateQrCode(req *types.GenerateQrCodeReq) (resp *types.GenerateQrCodeRes, err error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 生成场景 ID
	sceneBytes := make([]byte, 16)
	_, _ = rand.Read(sceneBytes)
	sceneID := hex.EncodeToString(sceneBytes)

	// 3. 创建二维码记录
	now := time.Now()
	qrCode := open_models.OpenQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0,                        // 等待扫码
		ExpiresAt: now.Add(5 * time.Minute), // 5分钟过期
	}

	if err := l.svcCtx.DB.Create(&qrCode).Error; err != nil {
		logx.Errorf("创建二维码记录失败: %v", err)
		return nil, errors.New("生成二维码失败")
	}

	// 4. 返回二维码信息
	return &types.GenerateQrCodeRes{
		SceneID:   sceneID,
		QrCodeURL: "https://api.beaver.im/qr/" + sceneID, // 实际应该生成真实的二维码图片 URL
		ExpireIn:  300,                                   // 5分钟
	}, nil
}
