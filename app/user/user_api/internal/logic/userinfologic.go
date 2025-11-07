package logic

import (
	"context"
	"fmt"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoRes, err error) {
	fmt.Println("获取用户的基础信息, UserID: %v,\n", req.UserID)

	// 直接从数据库查询，避免RPC调用自身服务
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "uuid = ?", req.UserID).Error
	if err != nil {
		fmt.Printf("[ERROR] 查询用户失败, UserID: %v, error: %v\n", req.UserID, err)
		return nil, err
	}

	resp = &types.UserInfoRes{
		UserID:   user.UUID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Abstract: user.Abstract,
		Phone:    user.Phone,
		Email:    user.Email,
		Gender:   user.Gender,
	}

	return resp, nil
}
