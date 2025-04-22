package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRemoveLogic {
	return &GroupMemberRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRemoveLogic) GroupMemberRemove(req *types.GroupMemberRemoveReq) (resp *types.GroupMemberRemoveRes, err error) {
	// 检查操作者权限
	var operator group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&operator, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("操作者不是群组成员")
	}
	if !(operator.Role == 1 || operator.Role == 2) {
		return nil, errors.New("没有权限移除成员")
	}

	// 检查要移除的成员
	var members []group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? and user_id in ?", req.GroupID, req.MemberIDs).Find(&members).Error
	if err != nil {
		return nil, errors.New("查询成员信息失败")
	}

	// 检查权限
	for _, member := range members {
		// 群主可以移除管理员和普通成员
		if operator.Role == 1 {
			if member.Role == 1 {
				return nil, errors.New("不能移除群主")
			}
		} else if operator.Role == 2 {
			// 管理员只能移除普通成员
			if member.Role != 3 {
				return nil, errors.New("管理员只能移除普通成员")
			}
		}
	}

	// 执行移除操作
	err = l.svcCtx.DB.Where("group_id = ? and user_id in ?", req.GroupID, req.MemberIDs).Delete(&group_models.GroupMemberModel{}).Error
	if err != nil {
		l.Logger.Errorf("移除成员失败: %v", err)
		return nil, errors.New("移除成员失败")
	}

	// 异步通知群成员
	defer func() {
		// 获取群成员列表
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(l.ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			l.Logger.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 通过ws推送给群成员
		for _, member := range response.Members {
			if member.UserID != req.UserID { // 不通知操作者自己
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberUpdate, req.UserID, member.UserID, map[string]interface{}{
					"groupId":   req.GroupID,
					"memberIds": req.MemberIDs,
					"operator":  req.UserID,
				})
			}
		}

		// 通知被移除的成员
		for _, memberID := range req.MemberIDs {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberUpdate, req.UserID, memberID, map[string]interface{}{
				"groupId":  req.GroupID,
				"operator": req.UserID,
				"memberId": memberID,
			})
		}
	}()

	return &types.GroupMemberRemoveRes{}, nil
}
