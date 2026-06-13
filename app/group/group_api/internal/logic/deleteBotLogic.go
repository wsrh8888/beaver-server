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
	"beaver/app/open/open_rpc/types/open_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type DeleteBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 删除机器人
func NewDeleteBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBotLogic {
	return &DeleteBotLogic{
		ctx:    ctx,
		logger: logger.New("delete_bot"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteBotLogic) DeleteBot(req *types.DeleteBotReq) (resp *types.DeleteBotRes, err error) {
	// 1. 从本地引用表查（通过 bot_id）
	var ref group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("bot_id = ?", req.BotID).First(&ref).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	// 2. 校验权限
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", ref.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可删除机器人")
	}

	// 3. 通过 open_rpc 获取 Bot ID，然后删除
	botInfoRes, err := l.svcCtx.OpenRpc.GetBotInfo(l.ctx, &open_rpc.GetBotInfoReq{
		BotId: ref.BotID,
	})
	if err != nil {
		return nil, errors.New("Open Bot 记录不存在")
	}

	if _, err = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{
		Id: botInfoRes.Id,
	}); err != nil {
		return nil, errors.New("删除失败")
	}

	// 4. 删本地引用表
	l.svcCtx.DB.Delete(&ref)

	// 5. 将机器人移出群（软删除）
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", ref.GroupID)
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ? AND user_id = ?", ref.GroupID, ref.BotID).
		Updates(map[string]interface{}{
			"status":  0,
			"version": memberVersion,
		})

	l.logger.Info(model.LogMsg{
		Text: "群机器人删除成功",
		Data: map[string]interface{}{
			"groupId": ref.GroupID,
			"userId":  req.UserID,
			"botId":   req.BotID,
		},
	})

	go l.notifyBotRemoved(ref.GroupID, req.UserID, memberVersion)

	return &types.DeleteBotRes{}, nil
}

func (l *DeleteBotLogic) notifyBotRemoved(groupID, operatorID string, memberVersion int64) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("异步通知群机器人删除时发生panic: %v", r)
		}
	}()

	ctx := context.Background()
	conversationID := "group_" + groupID

	_, err := l.svcCtx.ChatRpc.SendNotificationMessage(ctx, &chat_rpc.SendNotificationMessageReq{
		ConversationId: conversationID,
		MessageType:    8,
		Content:        fmt.Sprintf("%s 移除了群机器人", operatorID),
		RelatedUserId:  operatorID,
	})
	if err != nil {
		logx.Errorf("发送群机器人删除通知失败: groupId=%s, error=%v", groupID, err)
	}

	if memberVersion <= 0 {
		return
	}

	response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
		GroupID: groupID,
	})
	if err != nil {
		logx.Errorf("获取群成员列表失败: groupId=%s, error=%v", groupID, err)
		return
	}

	for _, member := range response.Members {
		payload := map[string]interface{}{
			"command":  wsCommandConst.GROUP_OPERATION,
			"type":     wsTypeConst.GroupMemberReceive,
			"senderId": operatorID,
			"targetId": member.UserID,
			"body": map[string]interface{}{
				"table": "group_members",
				"data": []map[string]interface{}{
					{
						"version": memberVersion,
						"groupId": groupID,
					},
				},
			},
			"conversationId": "",
		}
		if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logx.Errorf("推送群成员更新失败: groupId=%s, targetId=%s, error=%v", groupID, member.UserID, err)
		}
	}
}
