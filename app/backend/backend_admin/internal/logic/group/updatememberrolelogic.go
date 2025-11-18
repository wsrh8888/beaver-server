package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateMemberRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新群成员角色
func NewUpdateMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberRoleLogic {
	return &UpdateMemberRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMemberRoleLogic) UpdateMemberRole(req *types.UpdateMemberRoleReq) (resp *types.UpdateMemberRoleRes, err error) {
	// 检查群组成员是否存在
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("群组成员不存在: %d", req.Id)
			return nil, errors.New("群组成员不存在")
		}
		logx.Errorf("查询群组成员失败: %v", err)
		return nil, errors.New("查询群组成员失败")
	}

	// 验证角色值的有效性
	if req.Role < 1 || req.Role > 3 {
		return nil, errors.New("角色值无效，有效值为1-3")
	}

	// 更新成员角色
	err = l.svcCtx.DB.Model(&member).Update("role", req.Role).Error
	if err != nil {
		logx.Errorf("更新成员角色失败: %v", err)
		return nil, errors.New("更新成员角色失败")
	}

	return &types.UpdateMemberRoleRes{}, nil
}
