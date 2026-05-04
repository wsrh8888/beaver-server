package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.GetUserDetailReq) (resp *types.GetUserDetailRes, err error) {
	// 1. 查询用户信息
	var user user_models.UserModel
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 转换为响应格式
	return &types.GetUserDetailRes{
		User: types.UserDetail{
			UserID:   user.UserID,
			NickName: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
			Gender:   int(user.Gender),
			Status:   int(user.Status),
		},
	}, nil
}
