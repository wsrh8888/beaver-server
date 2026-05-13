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
	qrCode := open_models.OpenQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0, // 0-等待扫码
		ExpiresAt: expiresAt,
	}

	if err := l.svcCtx.DB.Create(&qrCode).Error; err != nil {
		logx.Errorf("创建扫码记录失败: err=%v", err)
		return nil, fmt.Errorf("服务内部异常")
	}

	// 6. 返回结果
	// 注意：QrCodeURL 字段用于兼容，实际前端应该使用 sceneId 自己生成二维码
	logx.Infof("生成扫码二维码成功: sceneId=%s, appId=%s", sceneID, req.AppID)

	return &types.GenerateQrCodeRes{
		QrCodeURL: "",      // 不使用，前端用 sceneId 生成二维码
		SceneID:   sceneID, // 前端用这个生成二维码
		ExpireIn:  300,     // 5分钟 = 300秒
	}, nil
}
