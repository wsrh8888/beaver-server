package oauth_public

import (
	"context"
	"fmt"
	"time"

	oauthmiddle "beaver/app/open/open_api/internal/middle/oauth"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrCodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckQrCodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrCodeStatusLogic {
	return &CheckQrCodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckQrCodeStatusLogic) CheckQrCodeStatus(req *types.CheckQrCodeStatusReq) (resp *types.CheckQrCodeStatusRes, err error) {
	var qrCode open_models.OpenOAuthQrCode
	if err := l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error; err != nil {
		logx.Errorf("扫码记录不存在: sceneId=%s, err=%v", req.SceneID, err)
		return nil, fmt.Errorf("二维码不存在或已过期")
	}

	if time.Now().After(qrCode.ExpiresAt) {
		l.svcCtx.DB.Model(&qrCode).Update("status", oauthmiddle.QrStatusExpired)
		return &types.CheckQrCodeStatusRes{Status: "expired"}, nil
	}

	status := oauthmiddle.QrStatusText(qrCode.Status)

	var code string
	if qrCode.Status == oauthmiddle.QrStatusConfirmed {
		if c, findErr := l.svcCtx.OAuth.FindConfirmedCode(req.SceneID, &qrCode); findErr == nil {
			code = c
		}
	}

	logx.Infof("查询扫码状态: sceneId=%s, status=%s", req.SceneID, status)

	return &types.CheckQrCodeStatusRes{
		Status: status,
		Code:   code,
	}, nil
}
