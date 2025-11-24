package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberAddLogic {
	return &GroupMemberAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberAddLogic) GroupMemberAdd(req *types.GroupMemberAddReq) (resp *types.GroupMemberAddRes, err error) {
	// 群成员邀请好友，isInvite 为true
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error

	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}

	// 检查哪些用户已经进群了，并分类处理
	var existingMembers []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&existingMembers, "group_id = ? and user_id in ?", req.GroupID, req.UserIds)

	// 构建已存在成员的映射（按userID）
	existingMemberMap := make(map[string]*group_models.GroupMemberModel)
	for i := range existingMembers {
		existingMemberMap[existingMembers[i].UserID] = &existingMembers[i]
	}

	// 检查是否有正常状态的成员（不允许重复添加）
	for _, existingMember := range existingMembers {
		if existingMember.Status == 1 {
			return nil, errors.New("已经有用户已经是群成员")
		}
	}

	// 分类处理：需要创建的新成员和需要更新的已存在成员
	var newMembers []group_models.GroupMemberModel
	var updateMembers []group_models.GroupMemberModel
	var lastVersion int64 // 记录最后一个版本号，用于返回

	for _, memberID := range req.UserIds {
		// 获取该群成员的版本号（按群独立递增）
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
		if memberVersion == -1 {
			l.Logger.Errorf("获取群成员版本号失败")
			return nil, errors.New("获取版本号失败")
		}
		lastVersion = memberVersion // 记录最后一个版本号

		existingMember, exists := existingMemberMap[memberID]
		if exists {
			// 成员已存在但状态不是正常（Status=2退出 或 3被踢），更新状态为正常
			updateMembers = append(updateMembers, group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   memberID,
				Role:     existingMember.Role, // 保持原有角色
				Status:   1,                   // 更新为正常状态
				JoinTime: time.Now(),          // 更新加入时间
				Version:  memberVersion,       // 更新版本号
			})
		} else {
			// 成员不存在，创建新记录
			newMembers = append(newMembers, group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   memberID,
				Role:     3, // 普通成员
				Status:   1, // 正常状态
				JoinTime: time.Now(),
				Version:  memberVersion,
			})
		}
	}

	// 创建新成员
	if len(newMembers) > 0 {
		err = l.svcCtx.DB.Create(&newMembers).Error
		if err != nil {
			l.Logger.Errorf("添加群成员失败: %v", err)
			return nil, errors.New("添加失败")
		}
	}

	// 更新已存在成员的状态
	if len(updateMembers) > 0 {
		for _, updateMember := range updateMembers {
			err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
				Where("group_id = ? AND user_id = ?", updateMember.GroupID, updateMember.UserID).
				Updates(map[string]interface{}{
					"status":    1,                    // 更新为正常状态
					"join_time": time.Now(),           // 更新加入时间
					"version":   updateMember.Version, // 更新版本号
				}).Error
			if err != nil {
				l.Logger.Errorf("更新群成员状态失败: %v", err)
				return nil, errors.New("更新成员状态失败")
			}
		}
	}

	// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel

	// 更新新成员的会话记录
	_, err = l.svcCtx.ChatRpc.BatchUpdateConversation(l.ctx, &chat_rpc.BatchUpdateConversationReq{
		UserIds:        req.UserIds,
		ConversationId: req.GroupID,
		LastMessage:    "",
	})
	if err != nil {
		logx.Errorf("Failed to update conversation: %v", err)
	}

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

		// 构建新加入成员的ID集合
		newMemberIds := make(map[string]bool)
		for _, newMemberID := range req.UserIds {
			newMemberIds[newMemberID] = true
		}

		// 获取群组版本号（用于通知群组信息变化）
		groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
		if groupVersion == -1 {
			l.Logger.Errorf("获取群组版本号失败")
			// 这里不影响主要功能，只是日志记录失败
		}

		// 预构建成员版本号映射表，避免重复查找
		memberVersionMap := make(map[string]int64)
		for _, member := range append(newMembers, updateMembers...) {
			memberVersionMap[member.UserID] = member.Version
		}

		// 1. 通知已在群的成员：group_members变化（成员列表增加了）
		// 构建新加入成员的数据列表
		var newMemberData []map[string]interface{}
		for _, newMemberID := range req.UserIds {
			newMemberVersion := memberVersionMap[newMemberID] // 直接从映射表获取

			newMemberData = append(newMemberData, map[string]interface{}{
				"version": newMemberVersion,
				"groupId": req.GroupID,
				"userId":  newMemberID,
			})
		}

		for _, member := range response.Members {
			if !newMemberIds[member.UserID] { // 不通知新加入的成员（已在群的所有成员都要收到通知，包括操作者自己）
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberReceive, req.UserID, member.UserID, map[string]interface{}{
					"tables": []map[string]interface{}{
						{
							"table": "group_members",
							"data":  newMemberData, // 推送所有新加入成员的信息列表
						},
					},
				}, "")
			}
		}

		// 2. 通知新加入的成员：group_members变化（他们成为了成员）+ groups变化（群基本信息）
		for _, newMemberID := range req.UserIds {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberReceive, req.UserID, newMemberID, map[string]interface{}{
				"tables": []map[string]interface{}{
					{
						"table": "groups",
						"data": []map[string]interface{}{
							{
								"version": groupVersion,
								"groupId": req.GroupID,
							},
						},
					},
					{
						"table": "group_members",
						"data":  newMemberData, // 推送所有新加入成员的信息列表
					},
				},
			}, "")
		}
	}()

	// 创建并返回响应
	resp = &types.GroupMemberAddRes{
		Version: lastVersion,
	}

	l.Logger.Infof("成功添加 %d 位成员到群组 %d", len(req.UserIds), req.GroupID)
	return resp, nil
}
