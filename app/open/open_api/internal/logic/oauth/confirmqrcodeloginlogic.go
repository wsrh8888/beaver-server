package oauth

import (
	"context"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	models "beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmQrCodeLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 确认扫码登录
func NewConfirmQrCodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrCodeLoginLogic {
	return &ConfirmQrCodeLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmQrCodeLoginLogic) ConfirmQrCodeLogin(req *types.ConfirmQrCodeLoginReq) (resp *types.ConfirmQrCodeLoginRes, err error) {
	// 1. 从中间件获取当前用户ID (需要 AuthMiddleware)
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, fmt.Errorf("未登录")
	}

	// 2. 查询扫码记录
	var qrCode models.OpenQrCode
	err = l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error
	if err != nil {
		return nil, err
	}

	// 3. 检查是否过期
	now := time.Now()
	if qrCode.ExpiresAt.Before(now) {
		return nil, fmt.Errorf("二维码已过期")
	}

	// 4. 检查状态
	if qrCode.Status != 0 && qrCode.Status != 1 {
		return nil, fmt.Errorf("二维码状态异常")
	}

	// 5. 更新状态为已确认,并绑定用户
	err = l.svcCtx.DB.Model(&qrCode).Updates(map[string]interface{}{
		"status":     2, // 已确认
		"user_id":    userID,
		"updated_at": now,
	}).Error
	if err != nil {
		return nil, err
	}

	logx.Infof("用户确认扫码登录: user_id=%s, scene_id=%s", userID, req.SceneID)

	return &types.ConfirmQrCodeLoginRes{
		Success: true,
	}, nil
}
