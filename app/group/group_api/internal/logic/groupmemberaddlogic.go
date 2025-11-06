package logic

import (
	"context"
	"errors"

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
	var memberVersion int64

	// 群成员邀请好友，isInvite 为true
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error

	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}

	// 去查一下哪些用户已经进群了。
	var memberList []group_models.GroupMemberModel

	l.svcCtx.DB.Find(&memberList, "group_id = ? and user_id in ?", req.GroupID, req.UserIds)

	if len(memberList) > 0 {
		return nil, errors.New("已经有用户已经是群成员")
	}

	// 清空memberList重新使用
	memberList = []group_models.GroupMemberModel{}

	for _, memberID := range req.UserIds {
		// 获取该群成员的版本号（按群独立递增）
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
		if memberVersion == -1 {
			l.Logger.Errorf("获取群成员版本号失败")
			return nil, errors.New("获取版本号失败")
		}

		memberList = append(memberList, group_models.GroupMemberModel{
			GroupID: req.GroupID,
			UserID:  memberID,
			Role:    3,
			Version: memberVersion,
		})
	}

	err = l.svcCtx.DB.Create(&memberList).Error

	if err != nil {
		l.Logger.Errorf("添加群成员失败: %v", err)
		return nil, errors.New("添加失败")
	}

	// 更新群组的成员版本号
	err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
		Where("group_id = ?", req.GroupID).
		Update("member_version", l.svcCtx.DB.Raw("member_version + 1")).Error
	if err != nil {
		l.Logger.Errorf("更新群组成员版本失败: %v", err)
		// 这里不返回错误，因为主要功能已经完成
	}

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

		// 通过ws推送给已存在的群成员（不通知操作者自己和新加入的成员）
		for _, member := range response.Members {
			if member.UserID != req.UserID && !newMemberIds[member.UserID] { // 不通知操作者自己和新加入的成员
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberUpdate, req.UserID, member.UserID, map[string]interface{}{
					"groupId":  req.GroupID,
					"type":     "add",
					"userIds":  req.UserIds,
					"operator": req.UserID,
				}, "")
			}
		}

		// 通知新加入的成员（需要获取完整的群组信息）
		for _, newMemberID := range req.UserIds {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberUpdate, req.UserID, newMemberID, map[string]interface{}{
				"groupId":      req.GroupID,
				"type":         "joined",
				"operator":     req.UserID,
				"needFullInfo": true, // 标记需要获取完整信息
			}, "")
		}
	}()

	// 创建并返回响应
	resp = &types.GroupMemberAddRes{
		Version: memberVersion,
	}

	l.Logger.Infof("成功添加 %d 位成员到群组 %d", len(req.UserIds), req.GroupID)
	return resp, nil
}
