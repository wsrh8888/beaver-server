package oauth

import (
	"context"
	"errors"
	"time"

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

// 查询扫码状态
func NewCheckQrCodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrCodeStatusLogic {
	return &CheckQrCodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckQrCodeStatusLogic) CheckQrCodeStatus(req *types.CheckQrCodeStatusReq) (resp *types.CheckQrCodeStatusRes, err error) {
	// 1. 查询二维码记录
	var qrCode open_models.OpenQrCode
	if err := l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error; err != nil {
		return nil, errors.New("二维码不存在")
	}

	// 2. 检查是否过期
	if time.Now().After(qrCode.ExpiresAt) {
		// 更新状态为已过期
		l.svcCtx.DB.Model(&qrCode).Update("status", 4)
		return &types.CheckQrCodeStatusRes{
			Status: "expired",
		}, nil
	}

	// 3. 映射状态码到字符串
	var statusStr string
	switch qrCode.Status {
	case 0:
		statusStr = "waiting"
	case 1:
		statusStr = "scanned"
	case 2:
		statusStr = "confirmed"
	case 3:
		statusStr = "cancelled"
	case 4:
		statusStr = "expired"
	default:
		statusStr = "unknown"
	}

	// 4. 返回状态
	return &types.CheckQrCodeStatusRes{
		Status: statusStr,
		UserId: qrCode.UserID,
	}, nil
}
