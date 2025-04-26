package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		CreatorID:  req.UserID,
		UUID:       utils.GenerateUUId(),
		Abstract:   "本群创建于" + time.Now().Format("2006-01-02") + "，欢迎大家加入",
		MaxMembers: 50,
		Avatar:     "cba984dd-bbe7-4ab4-80d7-b25654066f8d",
	}
	var groupUserList = []string{string(req.UserID)}
	if len(req.UserIdList) == 0 {
		return nil, errors.New("请选择用户")
	}

	for _, u := range req.UserIdList {
		groupUserList = append(groupUserList, u)
	}

	groupModel.Title = req.Name

	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Errorf("创建群失败: %v", err)
		return nil, errors.New("创建群失败")
	}

	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {

		memberMode := group_models.GroupMemberModel{
			GroupID: groupModel.UUID,
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

	// 异步通知群成员
	defer func() {
		// 获取群成员列表
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(l.ctx, &group_rpc.GetGroupMembersReq{
			GroupID: groupModel.UUID,
		})
		if err != nil {
			l.Logger.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 更新所有成员的会话记录
		allUserIDs := append([]string{req.UserID}, req.UserIdList...)
		fmt.Println("更新所有会话列表信息")
		_, err = l.svcCtx.ChatRpc.BatchUpdateConversation(l.ctx, &chat_rpc.BatchUpdateConversationReq{
			UserIds:        allUserIDs,
			ConversationId: groupModel.UUID,
			LastMessage:    "",
		})
		// 通过ws推送给群成员
		for _, member := range response.Members {
			fmt.Println("推送给群成员")
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.MessageGroupCreate, req.UserID, member.UserID, map[string]interface{}{
				"avatar":         groupModel.Avatar,
				"conversationId": groupModel.UUID,
				"update_at":      groupModel.CreatedAt.String(),
				"is_top":         false,
				"msg_preview":    "",
				"nickname":       groupModel.Title,
			}, groupModel.UUID)
		}
	}()

	return &types.GroupCreateRes{
		GroupID: groupModel.UUID,
	}, nil
}
