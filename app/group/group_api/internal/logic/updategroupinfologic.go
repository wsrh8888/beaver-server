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

type UpdateGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGroupInfoLogic) UpdateGroupInfo(req *types.UpdateGroupInfoReq) (resp *types.UpdateGroupInfoRes, err error) {
	// 检查操作者权限
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("操作者不是群组成员")
	}
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("没有权限更新群组信息")
	}

	// 构建更新字段
	updateFields := make(map[string]interface{})
	if req.Name != "" {
		updateFields["title"] = req.Name
	}
	if req.Avatar != "" {
		updateFields["avatar"] = req.Avatar
	}
	if req.Notice != "" {
		updateFields["notice"] = req.Notice
	}

	// 执行更新
	if len(updateFields) > 0 {
		err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
			Where("uuid = ?", req.GroupID).
			Updates(updateFields).Error
		if err != nil {
			l.Logger.Errorf("更新群组信息失败: %v", err)
			return nil, errors.New("更新群组信息失败")
		}
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
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupUpdate, req.UserID, member.UserID, map[string]interface{}{
					"groupId":  req.GroupID,
					"name":     req.Name,
					"avatar":   req.Avatar,
					"notice":   req.Notice,
					"joinType": req.JoinType,
				}, "")
			}
		}
	}()

	return &types.UpdateGroupInfoRes{}, nil
}
