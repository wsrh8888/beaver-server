package oauth

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type ScanQrCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewScanQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanQrCodeLogic {
	return &ScanQrCodeLogic{
		ctx:    ctx,
		logger: logger.New("scan_qrcode"),
		svcCtx: svcCtx,
	}
}

func (l *ScanQrCodeLogic) ScanQrCode(req *types.ScanQrCodeReq) (resp *types.ScanQrCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.SceneID == "" {
		return nil, errors.New("sceneId 不能为空")
	}

	if err := l.svcCtx.OAuth.MarkScanned(req.SceneID, req.UserID); err != nil {
		return nil, err
	}

	logx.Infof("扫码会话已标记 scanned: sceneId=%s, userId=%s", req.SceneID, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "OAuth扫码成功",
		Data: map[string]interface{}{
			"sceneId": req.SceneID,
			"userId":  req.UserID,
		},
	})

	return &types.ScanQrCodeRes{Success: true}, nil
}
