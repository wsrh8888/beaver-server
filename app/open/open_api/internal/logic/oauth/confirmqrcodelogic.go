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


type ConfirmQrCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewConfirmQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrCodeLogic {
	return &ConfirmQrCodeLogic{
		ctx:    ctx,
		logger: logger.New("confirm_qrcode"),
		svcCtx: svcCtx,
	}
}

func (l *ConfirmQrCodeLogic) ConfirmQrCode(req *types.ConfirmQrCodeReq) (resp *types.ConfirmQrCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.SceneID == "" {
		return nil, errors.New("sceneId 不能为空")
	}

	if err := l.svcCtx.OAuth.Confirm(req.SceneID, req.UserID); err != nil {
		return nil, err
	}

	logx.Infof("扫码授权已确认: sceneId=%s, userId=%s", req.SceneID, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "OAuth扫码确认成功",
		Data: map[string]interface{}{
			"sceneId": req.SceneID,
			"userId":  req.UserID,
		},
	})

	return &types.ConfirmQrCodeRes{Success: true}, nil
}
