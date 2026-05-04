// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	models "beaver/app/open/open_models"

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
	var app models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, errors.New("应用未启用")
	}

	// 3. 从 context 获取当前登录用户 ID
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("用户未登录")
	}

	// 4. 生成授权码
	codeBytes := make([]byte, 32)
	_, err = rand.Read(codeBytes)
	if err != nil {
		return nil, errors.New("生成授权码失败")
	}
	code := hex.EncodeToString(codeBytes)

	// 5. 保存授权码到数据库
	authCode := models.OpenAuthCode{
		Code:        code,
		AppID:       req.AppID,
		UserID:      userID,
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
		State:       req.State,
		ExpiresAt:   time.Now().Add(10 * time.Minute).Unix(), // 10分钟过期
		Used:        false,
	}

	err = l.svcCtx.DB.Create(&authCode).Error
	if err != nil {
		return nil, errors.New("保存授权码失败")
	}

	logx.Infof("生成授权码: app_id=%s, user_id=%s", req.AppID, userID)

	return &types.AuthorizeCodeRes{
		Code:  code,
		State: req.State,
	}, nil
}
