package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAppSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetAppSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAppSecretLogic {
	return &ResetAppSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAppSecretLogic) ResetAppSecret(req *types.ResetAppSecretReq) (resp *types.ResetAppSecretRes, err error) {
	// 生成新的 AppSecret（32字节随机字符串）
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}
	newSecret := hex.EncodeToString(secretBytes)

	// 更新数据库
	err = l.svcCtx.DB.Model(&open_models.OpenApp{}).
		Where("app_id = ?", req.AppID).
		Update("app_secret", newSecret).Error

	if err != nil {
		return nil, err
	}

	// 记录操作日志
	l.Infof("应用 %s 的密钥已重置，操作时间: %s", req.AppID, time.Now().Format("2006-01-02 15:04:05"))

	return &types.ResetAppSecretRes{
		AppID:     req.AppID,
		AppSecret: newSecret,
		Message:   "密钥重置成功，请妥善保管新密钥，旧密钥将立即失效",
	}, nil
}
