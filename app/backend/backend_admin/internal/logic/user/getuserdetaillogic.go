package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户详情
func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.GetUserDetailReq) (resp *types.GetUserDetailRes, err error) {
	var user user_models.UserModel

	err = l.svcCtx.DB.Where("uuid = ?", req.UserID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Logger.Errorf("用户不存在: %s", req.UserID)
			return nil, errors.New("用户不存在")
		}
		l.Logger.Errorf("查询用户详情失败: %v", err)
		return nil, errors.New("查询用户详情失败")
	}

	l.Logger.Infof("获取用户详情成功: userID=%s", req.UserID)
	return &types.GetUserDetailRes{
		Id:          user.UUID,
		Nickname:    user.NickName,
		FileName:    user.Avatar,
		Email:       user.Email,
		Abstract:    user.Abstract,
		Status:      int(user.Status),
		Source:      int(user.Source),
		LastLoginIP: "", // UserModel 没有 LastLoginIP 字段
		CreateTime:  user.CreatedAt.String(),
		UpdateTime:  user.UpdatedAt.String(),
	}, nil
}
