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

type GetAuthorizeCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取授权码
func NewGetAuthorizeCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorizeCodeLogic {
	return &GetAuthorizeCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuthorizeCodeLogic) GetAuthorizeCode(req *types.GetAuthorizeCodeReq) (resp *types.AuthorizeCodeRes, err error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 生成授权码
	codeBytes := make([]byte, 32)
	_, _ = rand.Read(codeBytes)
	code := hex.EncodeToString(codeBytes)

	// 3. 获取当前用户ID（从 context 中获取，需要通过中间件注入）
	userID, ok := l.ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录或登录已过期")
	}

	// 4. 保存授权码到数据库
	now := time.Now()
	authCode := open_models.OpenOAuthCode{
		Code:        code,
		AppID:       req.AppID,
		UserID:      userID,
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
		State:       req.State,
		ExpiresAt:   now.Add(10 * time.Minute).Unix(), // 10分钟过期
		Used:        false,
	}

	if err := l.svcCtx.DB.Create(&authCode).Error; err != nil {
		logx.Errorf("创建授权码失败: %v", err)
		return nil, errors.New("生成授权码失败")
	}

	// 5. 返回授权码
	return &types.AuthorizeCodeRes{
		Code:  code,
		State: req.State,
	}, nil
}
