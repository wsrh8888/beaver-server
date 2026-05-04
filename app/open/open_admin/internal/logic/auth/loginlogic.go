package auth

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	user_models "beaver/app/user/user_models"
	"beaver/utils/jwts"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRes, err error) {
	// 1. 根据手机号或邮箱查询用户
	var user user_models.UserModel
	err = l.svcCtx.DB.Where("phone = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		logx.Errorf("用户不存在: %s, error: %v", req.Username, err)
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 验证密码
	if !pwd.CheckPad(user.Password, req.Password) {
		logx.Errorf("密码错误: user_id=%s", user.UserID)
		return nil, errors.New("用户名或密码错误")
	}
	
	// 3. 检查是否是已认证的开发者
	var developerCount int64
	err = l.svcCtx.DB.Table("open_developers").
		Where("user_id = ? AND status = ?", user.UserID, 1).
		Count(&developerCount).Error
	
	if err != nil || developerCount == 0 {
		return nil, errors.New("您还不是认证开发者,请先申请开发者资质")
	}
	
	// 4. 生成 JWT Token (与 IM 登录保持一致)
	secretKey := l.svcCtx.Config.Auth.AccessSecret
	expireHours := l.svcCtx.Config.Auth.AccessExpire / 3600 // 转换为小时
	if expireHours == 0 {
		expireHours = 12 // 默认 12 小时
	}
	
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   user.UserID,
		NickName: user.NickName,
	}, secretKey, int(expireHours))
	
	if err != nil {
		logx.Errorf("生成 token 失败: %v", err)
		return nil, errors.New("服务内部异常")
	}
	
	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour).UnixMilli()
	
	logx.Infof("开放平台登录成功: user_id=%s, nick_name=%s", user.UserID, user.NickName)
	
	return &types.LoginRes{
		Token:    token,
		UserID:   user.UserID,
		NickName: user.NickName,
		ExpireAt: expireAt,
	}, nil
}
