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
	err = l.svcCtx.DB.Where("group_id = ? and user_id in ?", req.GroupID, req.UserIds).Find(&members).Error
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

	// 获取该群成员的版本号（按群独立递增）
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
	if memberVersion == -1 {
		l.Logger.Errorf("获取群成员版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	// 批量更新成员状态为被踢（Status = 3）
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id in ?", req.GroupID, req.UserIds).
		Updates(map[string]interface{}{
			"status":  3, // 3被踢
			"version": memberVersion,
		}).Error
	if err != nil {
		l.Logger.Errorf("移除成员失败: %v", err)
		return nil, errors.New("移除成员失败")
	}

	// 注意：群成员移除时，groups 表信息没有变化，不需要更新版本号

	// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel

	// 异步通知群成员
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			l.Logger.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 构建被移除成员的数据列表
		var removedMemberData []map[string]interface{}
		for _, memberID := range req.UserIds {
			removedMemberData = append(removedMemberData, map[string]interface{}{
				"version": memberVersion,
				"groupId": req.GroupID,
				"userId":  memberID,
			})
		}

		// 通知所有人（当前群成员 + 被移除的成员）：group_members变化
		// 前端收到通知后，自己判断是否还在这群里，如果不在了就删除群信息

		// 去重，避免重复通知
		memberMap := make(map[string]bool)
		var allMembers []string

		// 添加当前群成员
		for _, member := range response.Members {
			if !memberMap[member.UserID] {
				memberMap[member.UserID] = true
				allMembers = append(allMembers, member.UserID)
			}
		}

		// 添加被移除的成员
		for _, memberID := range req.UserIds {
			if !memberMap[memberID] {
				memberMap[memberID] = true
				allMembers = append(allMembers, memberID)
			}
		}

		// 统一给所有人发送相同的消息
		for _, memberID := range allMembers {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberReceive, req.UserID, memberID, map[string]interface{}{
				"tables": []map[string]interface{}{
					{
						"table": "group_members",
						"data":  removedMemberData, // 推送所有被移除成员的信息列表
					},
				},
			}, "")
		}
	}()

	return &types.GroupMemberRemoveRes{
		Version: memberVersion,
	}, nil
}
