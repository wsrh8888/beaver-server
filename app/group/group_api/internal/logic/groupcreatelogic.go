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

	//先生成群组ID
	groupID := utils.GenerateUUId()

	// 获取该群组的版本号（每个群独立递增）
	groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", groupID)
	if groupVersion == -1 {
		logx.Errorf("获取群组版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	var groupModel = group_models.GroupModel{
		CreatorID: req.UserID,
		GroupID:   groupID,
		Version:   groupVersion, // 该群的独立版本
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
		// 获取该群成员的版本号（按群独立递增）
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", groupID)

		memberMode := group_models.GroupMemberModel{
			GroupID:  groupID,
			UserID:   u,
			Role:     3,
			JoinTime: time.Now(),
			Version:  memberVersion, // 该群成员的独立版本
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

	// 创建成员变更日志
	var changeLogs []group_models.GroupMemberChangeLogModel
	for _, member := range members {
		// 获取全局递增的变更日志版本号
		logVersion := l.svcCtx.VersionGen.GetNextVersion("group_member_logs", "", "")
		if logVersion == -1 {
			logx.Errorf("获取变更日志版本号失败，用户ID: %s", member.UserID)
			return nil, errors.New("获取版本号失败")
		}

		changeLog := group_models.GroupMemberChangeLogModel{
			GroupID:    member.GroupID,
			UserID:     member.UserID,
			ChangeType: "join",
			OperatedBy: req.UserID, // 创建者操作
			ChangeTime: member.JoinTime,
			Version:    logVersion,
		}
		changeLogs = append(changeLogs, changeLog)
	}

	err = l.svcCtx.DB.Create(&changeLogs).Error
	if err != nil {
		logx.Errorf("创建群成员变更日志失败: %v", err)
		// 这里不返回错误，因为主要功能已经完成，只是日志记录失败
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
		Version: groupVersion,
	}, nil
}
