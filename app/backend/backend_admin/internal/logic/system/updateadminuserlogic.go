package system

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminUserLogic {
	return &UpdateAdminUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAdminUserLogic) UpdateAdminUser(req *types.UpdateAdminUserReq) (resp *types.UpdateAdminUserRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	var adminUser backend_models.AdminUser
	if err = l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&adminUser).Error; err != nil {
		return nil, errors.New("管理员不存在")
	}

	updates := map[string]interface{}{}
	if req.NickName != "" {
		updates["nick_name"] = req.NickName
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.Password != "" {
		updates["password"] = pwd.HahPwd(req.Password)
	}
	if len(updates) > 0 {
		if err = l.svcCtx.DB.Model(&adminUser).Updates(updates).Error; err != nil {
			l.Errorf("更新管理员失败: %v", err)
			return nil, err
		}
	}

	if req.AuthorityIds != nil {
		if err = l.svcCtx.DB.Where("user_id = ?", req.UserID).
			Delete(&backend_models.AdminSystemAuthorityUser{}).Error; err != nil {
			return nil, err
		}
		if len(req.AuthorityIds) > 0 {
			rows := make([]backend_models.AdminSystemAuthorityUser, 0, len(req.AuthorityIds))
			for _, aid := range req.AuthorityIds {
				rows = append(rows, backend_models.AdminSystemAuthorityUser{
					UserID:      req.UserID,
					AuthorityID: aid,
				})
			}
			if err = l.svcCtx.DB.Create(&rows).Error; err != nil {
				l.Errorf("更新管理员角色失败: %v", err)
				return nil, err
			}
		}
	}
	return &types.UpdateAdminUserRes{}, nil
}
