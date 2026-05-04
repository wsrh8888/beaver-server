package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetH5AuthCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// H5 免登获取 authCode（需在 WebView 环境中调用）
func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (resp *types.GetH5AuthCodeRes, err error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 获取当前用户ID（从 context 中获取，H5 WebView 环境中应该已经登录）
	userID, ok := l.ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录或登录已过期")
	}

	// 3. 生成 H5 AuthCode
	codeBytes := make([]byte, 32)
	_, _ = rand.Read(codeBytes)
	code := hex.EncodeToString(codeBytes)

	// 4. 保存 H5 AuthCode（5分钟过期）
	now := time.Now()
	h5AuthCode := open_models.OpenH5AuthCode{
		Code:      code,
		AppID:     req.AppID,
		UserID:    userID,
		ExpiresAt: now.Add(5 * time.Minute).Unix(),
		CreatedAt: now.Unix(),
	}

	if err := l.svcCtx.DB.Create(&h5AuthCode).Error; err != nil {
		logx.Errorf("创建 H5 AuthCode 失败: %v", err)
		return nil, errors.New("生成授权码失败")
	}

	return &types.GetH5AuthCodeRes{
		AuthCode: code,
	}, nil
}
