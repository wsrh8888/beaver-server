package oauth

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScanQrCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScanQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScanQrCodeLogic {
	return &ScanQrCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
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

	return &types.ScanQrCodeRes{Success: true}, nil
}
