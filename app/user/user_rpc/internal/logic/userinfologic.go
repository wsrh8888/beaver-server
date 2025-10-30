package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user_rpc.UserInfoReq) (*user_rpc.UserInfoRes, error) {
	var user user_models.UserModel

	err := l.svcCtx.DB.Take(&user, "uuid = ?", in.UserID).Error

	if err != nil {
		logx.Errorf("查询用户失败: %s", err.Error())
		return nil, errors.New("用户不存在")
	}

	// 创建安全的用户信息结构，不包含敏感信息
	safeUserInfo := map[string]interface{}{
		"uuid":      user.UUID,
		"nickName":  user.NickName,
		"email":     user.Email,
		"phone":     user.Phone,
		"avatar":    user.Avatar,
		"abstract":  user.Abstract,
		"status":    user.Status,
		"source":    user.Source,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		"gender":    user.Gender,
	}

	byteData, _ := json.Marshal(safeUserInfo)

	return &user_rpc.UserInfoRes{
		Data: byteData,
	}, nil
}
