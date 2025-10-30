package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	utils "beaver/utils/rand"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupCreateLogic) GroupCreate(req *types.GroupCreateReq) (resp *types.GroupCreateRes, err error) {

	var groupModel = group_models.GroupModel{
		CreatorID: req.UserID,
		GroupID:   utils.GenerateUUId(),
	}
	var groupUserList = []string{req.UserID}
	if len(req.UserIdList) == 0 {
		return nil, errors.New("请选择用户")
	}

	for _, u := range req.UserIdList {
		groupUserList = append(groupUserList, u)
	}

	groupModel.Title = req.Title

	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Errorf("创建群失败: %v", err)
		return nil, errors.New("创建群失败")
	}

	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {

		memberMode := group_models.GroupMemberModel{
			GroupID: groupModel.GroupID,
			UserID:  u,
			Role:    3,
		}
		if i == 0 {
			memberMode.Role = 1
		}
		members = append(members, memberMode)
	}

	err = l.svcCtx.DB.Create(&members).Error
	if err != nil {
		logx.Errorf("创建群成员失败: %v", err)
		return nil, errors.New("创建群成员失败")
	}

	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: groupModel.GroupID,
		})
		if err != nil {
			l.Logger.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 更新所有成员的会话记录
		allUserIDs := append([]string{req.UserID}, req.UserIdList...)
		fmt.Println("更新所有会话列表信息")
		_, err = l.svcCtx.ChatRpc.BatchUpdateConversation(ctx, &chat_rpc.BatchUpdateConversationReq{
			UserIds:        allUserIDs,
			ConversationId: groupModel.GroupID,
			LastMessage:    "",
		})

		// 通过ws推送给群成员
		for _, member := range response.Members {
			fmt.Println("推送给群成员")
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.MessageGroupCreate, req.UserID, member.UserID, map[string]interface{}{
				"avatar":         groupModel.Avatar,
				"conversationId": groupModel.GroupID,
				"update_at":      groupModel.CreatedAt.String(),
				"is_top":         false,
				"msg_preview":    "",
				"nickname":       groupModel.Title,
			}, groupModel.GroupID)
		}
	}()

	return &types.GroupCreateRes{
		GroupID: groupModel.GroupID,
	}, nil
}
