package oauth

import (
	"context"
	"fmt"
	"time"

	"beaver-server/app/open/open_models"
	"beaver-server/app/user/user_models"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByH5CodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoByH5CodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByH5CodeLogic {
	return &GetUserInfoByH5CodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserInfoByH5Code 用 H5 authCode 换取用户信息
func (l *GetUserInfoByH5CodeLogic) GetUserInfoByH5Code(req *types.GetUserInfoByH5CodeReq) (*types.GetUserInfoByH5CodeRes, error) {
	// 1. 验证应用凭证
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ?", req.AppID, req.AppSecret).First(&app).Error; err != nil {
		return nil, fmt.Errorf("应用凭证错误")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用已禁用")
	}

	// 3. 查询 authCode
	var h5AuthCode open_models.OpenH5AuthCode
	if err := l.svcCtx.DB.Where("code = ? AND app_id = ?", req.AuthCode, req.AppID).First(&h5AuthCode).Error; err != nil {
		return nil, fmt.Errorf("授权码无效")
	}

	// 4. 检查是否过期
	if time.Now().Unix() > h5AuthCode.ExpiresAt {
		return nil, fmt.Errorf("授权码已过期")
	}

	// 5. 删除已使用的 authCode（一次性使用）
	l.svcCtx.DB.Delete(&h5AuthCode)

	// 6. 获取用户信息
	var user user_models.User
	if err := l.svcCtx.DB.Where("id = ?", h5AuthCode.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return &types.GetUserInfoByH5CodeRes{
		UserID:   user.ID,
		NickName: user.Nickname,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}
