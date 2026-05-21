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
	var qrCode open_models.OpenOAuthQrCode
	if err := l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error; err != nil {
		logx.Errorf("扫码记录不存在: sceneId=%s, err=%v", req.SceneID, err)
		return nil, fmt.Errorf("二维码不存在或已过期")
	}

	// 2. 检查是否过期
	if time.Now().After(qrCode.ExpiresAt) {
		// 更新状态为已过期
		l.svcCtx.DB.Model(&qrCode).Update("status", 4)
		return &types.CheckQrCodeStatusRes{
			Status: "expired",
		}, nil
	}

	// 3. 根据状态返回结果
	var status string
	switch qrCode.Status {
	case 0:
		status = "waiting" // 等待扫码
	case 1:
		status = "scanned" // 已扫码，待确认
	case 2:
		status = "confirmed" // 已确认
	case 3:
		status = "cancelled" // 已取消
	case 4:
		status = "expired" // 已过期
	default:
		status = "waiting"
	}

	// 4. 如果已确认，返回用户ID
	var userId string
	if qrCode.Status == 2 {
		userId = qrCode.UserID
	}

	logx.Infof("查询扫码状态: sceneId=%s, status=%s", req.SceneID, status)

	return &types.CheckQrCodeStatusRes{
		Status: status,
		UserId: userId,
	}, nil
}
