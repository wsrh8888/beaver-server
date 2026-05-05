package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAppSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置应用密钥
func NewResetAppSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAppSecretLogic {
	return &ResetAppSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAppSecretLogic) ResetAppSecret(req *types.ResetAppSecretReq) (resp *types.ResetAppSecretRes, err error) {
	// 1. 从 header 获取当前用户 ID
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, errors.New("未登录")
	}

	// 2. 生成新密钥
	newSecret := uuid.New().String() + uuid.New().String()

	// 3. 更新密钥
	result := l.svcCtx.DB.Model(&open_models.OpenApp{}).
		Where("app_id = ? AND owner_user_id = ?", req.AppID, userID).
		Update("app_secret", newSecret)

	if result.Error != nil {
		logx.Errorf("重置密钥失败: %v", result.Error)
		return nil, errors.New("重置失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("应用不存在或无权限")
	}

	logx.Infof("应用密钥重置成功: app_id=%s", req.AppID)

	return &types.ResetAppSecretRes{
		AppSecret: newSecret,
	}, nil
}
