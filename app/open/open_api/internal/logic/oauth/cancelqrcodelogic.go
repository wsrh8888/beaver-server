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


type CancelQrCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewCancelQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelQrCodeLogic {
	return &CancelQrCodeLogic{
		ctx:    ctx,
		logger: logger.New("cancel_qrcode"),
		svcCtx: svcCtx,
	}
}

func (l *CancelQrCodeLogic) CancelQrCode(req *types.CancelQrCodeReq) (resp *types.CancelQrCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.SceneID == "" {
		return nil, errors.New("sceneId 不能为空")
	}

	if err := l.svcCtx.OAuth.Cancel(req.SceneID, req.UserID); err != nil {
		return nil, err
	}

	logx.Infof("扫码授权已取消: sceneId=%s, userId=%s", req.SceneID, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "OAuth扫码取消成功",
		Data: map[string]interface{}{
			"sceneId": req.SceneID,
			"userId":  req.UserID,
		},
	})

	return &types.CancelQrCodeRes{Success: true}, nil
}
