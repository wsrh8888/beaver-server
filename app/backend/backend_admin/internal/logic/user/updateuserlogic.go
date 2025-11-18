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
	// 检查用户是否存在
	var user user_models.UserModel
	err = l.svcCtx.DB.Where("uuid = ?", req.UserID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			l.Logger.Errorf("用户不存在: %s", req.UserID)
			return nil, errors.New("用户不存在")
		}
		l.Logger.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 如果要更新邮箱，检查是否重复
	if req.Email != nil && *req.Email != user.Email {
		var existUser user_models.UserModel
		err = l.svcCtx.DB.Where("email = ? AND uuid != ?", *req.Email, req.UserID).First(&existUser).Error
		if err == nil {
			l.Logger.Errorf("邮箱已存在: %s", *req.Email)
			return nil, errors.New("邮箱已存在")
		}
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.Nickname != nil {
		updates["nick_name"] = *req.Nickname
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.FileName != nil {
		updates["file_name"] = *req.FileName
	}
	if req.Abstract != nil {
		updates["abstract"] = *req.Abstract
	}
	if req.Status != nil {
		updates["status"] = int8(*req.Status)
	}

	// 执行更新
	if len(updates) > 0 {
		err = l.svcCtx.DB.Model(&user).Updates(updates).Error
		if err != nil {
			l.Logger.Errorf("更新用户失败: %v", err)
			return nil, errors.New("更新用户失败")
		}
		l.Logger.Infof("更新用户成功: userID=%s, updates=%v", req.UserID, updates)
	} else {
		l.Logger.Infof("用户信息无变化: userID=%s", req.UserID)
	}

	return &types.UpdateUserRes{}, nil
}
