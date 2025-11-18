package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除群成员
func NewRemoveGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMemberLogic {
	return &RemoveGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveGroupMemberLogic) RemoveGroupMember(req *types.RemoveGroupMemberReq) (resp *types.RemoveGroupMemberRes, err error) {
	if req.GroupId == "" {
		return nil, errors.New("群组ID不能为空")
	}

	if len(req.MemberIds) == 0 {
		return nil, errors.New("成员ID列表不能为空")
	}

	// 批量删除群组成员
	err = l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", req.GroupId, req.MemberIds).Delete(&group_models.GroupMemberModel{}).Error
	if err != nil {
		logx.Errorf("移除群组成员失败: %v", err)
		return nil, errors.New("移除群组成员失败")
	}

	return &types.RemoveGroupMemberRes{}, nil
}
