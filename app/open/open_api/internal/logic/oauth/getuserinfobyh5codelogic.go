package oauth

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByH5CodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用 H5 authCode 换取用户信息
func NewGetUserInfoByH5CodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByH5CodeLogic {
	return &GetUserInfoByH5CodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByH5CodeLogic) GetUserInfoByH5Code(req *types.GetUserInfoByH5CodeReq) (resp *types.GetUserInfoByH5CodeRes, err error) {
	// 1. 验证应用凭证
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用 ID 或密钥错误")
	}

	// 2. 查询 H5 AuthCode
	var h5AuthCode open_models.OpenH5AuthCode
	if err := l.svcCtx.DB.Where("code = ? AND app_id = ?", req.AuthCode, req.AppID).First(&h5AuthCode).Error; err != nil {
		return nil, errors.New("授权码无效")
	}

	// 3. 检查是否过期
	if time.Now().Unix() > h5AuthCode.ExpiresAt {
		return nil, errors.New("授权码已过期")
	}

	// 4. 查询用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", h5AuthCode.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 5. 删除已使用的授权码
	l.svcCtx.DB.Delete(&h5AuthCode)

	// 6. 返回用户信息
	return &types.GetUserInfoByH5CodeRes{
		UserID:   user.UserID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Email:    user.Email,
	}, nil
}
