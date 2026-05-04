package oauth

import (
	"context"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	models "beaver/app/open/open_models"

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
	// 1. 查询扫码记录
	var qrCode models.OpenQrCode
	err = l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error
	if err != nil {
		return nil, err
	}

	// 2. 检查是否过期
	now := time.Now()
	if qrCode.ExpiresAt.Before(now) {
		// 更新状态为已过期
		l.svcCtx.DB.Model(&qrCode).Update("status", 4)
		
		return &types.CheckQrCodeStatusRes{
			Status: "expired",
		}, nil
	}

	// 3. 返回当前状态
	statusMap := map[int]string{
		0: "waiting",
		1: "scanned",
		2: "confirmed",
		3: "cancelled",
		4: "expired",
	}

	res := &types.CheckQrCodeStatusRes{
		Status: statusMap[qrCode.Status],
	}

	// 如果已确认,返回用户信息
	if qrCode.Status == 2 && qrCode.UserID != "" {
		res.UserId = qrCode.UserID
	}

	return res, nil
}
