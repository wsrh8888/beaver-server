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
	// 1. 查询二维码记录
	var qrCode open_models.OpenQrCode
	if err := l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error; err != nil {
		return nil, errors.New("二维码不存在")
	}

	// 2. 检查是否过期
	if time.Now().After(qrCode.ExpiresAt) {
		l.svcCtx.DB.Model(&qrCode).Update("status", 4)
		return nil, errors.New("二维码已过期")
	}

	// 3. 检查状态（只能是等待扫码或已扫码状态）
	if qrCode.Status != 0 && qrCode.Status != 1 {
		return nil, errors.New("二维码状态不正确")
	}

	// 4. 获取当前用户ID（从 context 中获取）
	userID, ok := l.ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录或登录已过期")
	}

	// 5. 更新二维码状态为已确认
	if err := l.svcCtx.DB.Model(&qrCode).Updates(map[string]interface{}{
		"status":  2, // 已确认
		"user_id": userID,
	}).Error; err != nil {
		logx.Errorf("更新二维码状态失败: %v", err)
		return nil, errors.New("确认失败")
	}

	return &types.ConfirmQrCodeLoginRes{
		Success: true,
	}, nil
}
