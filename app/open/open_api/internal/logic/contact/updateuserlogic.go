package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.UpdateUserRes, err error) {
	// 1. 查询用户是否存在
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 更新用户信息
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nick_name"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) == 0 {
		return nil, errors.New("没有需要更新的字段")
	}

	if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
		logx.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	return &types.UpdateUserRes{
		Success: true,
	}, nil
}
