// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package oauth_public

import (
	"context"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/logic/oauthutil"
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
	// 1. 验证 userId
	if req.UserID == "" {
		logx.Error("确认扫码登录失败：userId为空")
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 2. 查询扫码记录
	var qrCode open_models.OpenOAuthQrCode
	if err := l.svcCtx.DB.Where("scene_id = ?", req.SceneID).First(&qrCode).Error; err != nil {
		logx.Errorf("扫码记录不存在: sceneId=%s, err=%v", req.SceneID, err)
		return nil, fmt.Errorf("二维码不存在或已过期")
	}

	// 3. 检查是否过期
	if time.Now().After(qrCode.ExpiresAt) {
		l.svcCtx.DB.Model(&qrCode).Update("status", 4)
		return nil, fmt.Errorf("二维码已过期")
	}

	// 4. 检查状态（只有等待扫码或已扫码状态才能确认）
	if qrCode.Status != 0 && qrCode.Status != 1 {
		return nil, fmt.Errorf("二维码状态不正确")
	}

	// 5. 更新扫码记录：设置用户ID和状态为已确认
	if err := l.svcCtx.DB.Model(&qrCode).Updates(map[string]interface{}{
		"user_id": req.UserID,
		"status":  2,
	}).Error; err != nil {
		logx.Errorf("更新扫码记录失败: err=%v", err)
		return nil, fmt.Errorf("服务内部异常")
	}

	if _, _, err := oauthutil.CreateOAuthCode(l.svcCtx.DB, qrCode.AppID, req.UserID, "pc_scan"); err != nil {
		logx.Errorf("生成扫码 authCode 失败: sceneId=%s, err=%v", req.SceneID, err)
		return nil, fmt.Errorf("生成授权码失败")
	}

	logx.Infof("确认扫码登录成功: sceneId=%s, userId=%s", req.SceneID, req.UserID)

	return &types.ConfirmQrCodeLoginRes{
		Success: true,
	}, nil
}
