package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	models "beaver/app/open/open_models"

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
	var app models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		return nil, err
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	// 3. 生成场景ID
	sceneBytes := make([]byte, 16)
	_, err = rand.Read(sceneBytes)
	if err != nil {
		return nil, err
	}
	sceneID := hex.EncodeToString(sceneBytes)

	// 4. 生成二维码URL (这里需要实现二维码生成逻辑)
	// 实际应该调用二维码生成服务或第三方API
	qrCodeURL := "https://open.beaver.im/scan/" + sceneID

	// 5. 保存扫码记录到数据库
	now := time.Now()
	qrCodeRecord := models.OpenQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0, // 0-等待扫码
		ExpiresAt: now.Add(5 * time.Minute),
	}

	err = l.svcCtx.DB.Create(&qrCodeRecord).Error
	if err != nil {
		return nil, err
	}

	logx.Infof("生成扫码二维码: app_id=%s, scene_id=%s", req.AppID, sceneID)

	return &types.GenerateQrCodeRes{
		QrCodeURL: qrCodeURL,
		SceneID:   sceneID,
		ExpireIn:  300, // 5分钟
	}, nil
}
