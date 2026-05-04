package user

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户信息
func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
	// 1. 查询用户列表
	var users []user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id IN ?", req.UserIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var userList []types.UserInfo
	for _, user := range users {
		userList = append(userList, types.UserInfo{
			UserID:   user.UserID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
		})
	}

	return &types.GetUserListRes{
		Users: userList,
	}, nil
}
