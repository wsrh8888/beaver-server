package oauth_public

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type PasswordLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码登录
func NewPasswordLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordLoginLogic {
	return &PasswordLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordLoginLogic) PasswordLogin(req *types.PasswordLoginReq) (resp *types.PasswordLoginRes, err error) {
	// 1. 验证 appId 是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	// 3. 根据账号查询用户（支持邮箱或工号）
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "email = ? OR user_id = ?", req.Account, req.Account).Error
	if err != nil {
		logx.Errorf("用户不存在: account=%s, err=%v", req.Account, err)
		return nil, errors.New("账号或密码错误")
	}

	// 4. 验证密码
	if !pwd.CheckPad(user.Password, req.Password) {
		logx.Errorf("密码错误: userId=%s", user.UserID)
		return nil, errors.New("账号或密码错误")
	}

	// 5. 生成 access_token 和 refresh_token
	accessToken := util.NewV4().String()
	refreshToken := util.NewV4().String()

	// 6. 设置过期时间（access_token 2小时）
	accessExpiresIn := int64(7200) // 2小时

	// 7. 存储 token 到数据库
	now := time.Now()
	tokenRecord := open_models.OpenAccessToken{
		AppID:        req.AppID,
		UserID:       user.UserID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(time.Duration(accessExpiresIn) * time.Second).Unix(),
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		logx.Errorf("存储 token 失败: err=%v", err)
		return nil, fmt.Errorf("服务内部异常")
	}

	logx.Infof("账号密码登录成功: userId=%s, appId=%s, account=%s", user.UserID, req.AppID, req.Account)

	return &types.PasswordLoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExpiresIn,
		UserID:       user.UserID,
		NickName:     user.NickName,
		Avatar:       user.Avatar,
	}, nil
}
