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
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	const expireIn int64 = 300
	sceneID := util.NewV4().String()
	expiresAt := time.Now().Add(time.Duration(expireIn) * time.Second)

	qrCode := open_models.OpenOAuthQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0,
		ExpiresAt: expiresAt,
	}
	if err := l.svcCtx.DB.Create(&qrCode).Error; err != nil {
		logx.Errorf("创建扫码记录失败: err=%v", err)
		return nil, fmt.Errorf("服务内部异常")
	}

	logx.Infof("生成扫码会话成功: sceneId=%s, appId=%s", sceneID, req.AppID)

	return &types.GenerateQrCodeRes{
		SceneID:  sceneID,
		ExpireIn: expireIn,
	}, nil
}
