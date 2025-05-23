package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

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

	// 去查一下哪些用户已经进群了。
	var memberList []group_models.GroupMemberModel

	l.svcCtx.DB.Find(&memberList, "group_id = ? and user_id in ?", req.GroupID, req.UserIds)

	if len(memberList) > 0 {
		return nil, errors.New("已经有用户已经是群成员")
	}

	// 清空memberList重新使用
	memberList = []group_models.GroupMemberModel{}

	for _, memberID := range req.UserIds {
		memberList = append(memberList, group_models.GroupMemberModel{
			GroupID: req.GroupID,
			UserID:  memberID,
			Role:    3,
		})
	}

	err = l.svcCtx.DB.Create(&memberList).Error

	if err != nil {
		l.Logger.Errorf("添加群成员失败: %v", err)
		return nil, errors.New("添加失败")
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

	// 创建并返回响应
	resp = &types.GroupMemberAddRes{}

	l.Logger.Infof("成功添加 %d 位成员到群组 %d", len(req.UserIds), req.GroupID)
	return resp, nil
}
