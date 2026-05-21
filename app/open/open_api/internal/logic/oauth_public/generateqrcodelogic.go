// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package oauth_public

import (
	"context"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	util "beaver/utils/uuid"

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
	// 1. 验证 appId 是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	// 3. 生成 SceneID
	sceneID := util.NewV4().String()

	// 4. 设置过期时间（5分钟）
	now := time.Now()
	expiresAt := now.Add(5 * time.Minute)

	// 5. 创建扫码记录
	qrCode := open_models.OpenOAuthQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0, // 0-等待扫码
		ExpiresAt: expiresAt,
	}

	if err := l.svcCtx.DB.Create(&qrCode).Error; err != nil {
		logx.Errorf("创建扫码记录失败: err=%v", err)
		return nil, fmt.Errorf("服务内部异常")
	}

	// 6. 生成短链接（TODO: 实现短链接服务）
	// 暂时使用完整 URL
	oauthBaseUrl := l.svcCtx.Config.OAuth.BaseUrl
	qrCodeURL := fmt.Sprintf("%s/scan?sceneId=%s", oauthBaseUrl, sceneID)

	// 7. 返回结果
	logx.Infof("生成扫码二维码成功: sceneId=%s, appId=%s, qrCodeUrl=%s", sceneID, req.AppID, qrCodeURL)

	return &types.GenerateQrCodeRes{
		QrCodeURL: qrCodeURL, // 二维码 URL（短链接或完整 URL）
		SceneID:   sceneID,
		ExpireIn:  300, // 5分钟 = 300秒
	}, nil
}
