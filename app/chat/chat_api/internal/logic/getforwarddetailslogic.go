package logic

import (
	"context"
	"encoding/json"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetForwardDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetForwardDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetForwardDetailsLogic {
	return &GetForwardDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetForwardDetailsLogic) GetForwardDetails(req *types.GetForwardDetailsReq) (resp *types.GetForwardDetailsRes, err error) {
	var detail chat_models.ChatForward
	err = l.svcCtx.DB.Where("record_id = ?", req.RecordID).First(&detail).Error
	if err != nil {
		l.Logger.Errorf("获取合并转发详情失败: %v", err)
		return nil, err
	}

	// 收集需要查询用户信息的UserID列表
	var userIds []string
	userIdSet := make(map[string]bool)
	for _, m := range detail.Content {
		if m.SendUserID != nil && *m.SendUserID != "" {
			if !userIdSet[*m.SendUserID] {
				userIds = append(userIds, *m.SendUserID)
				userIdSet[*m.SendUserID] = true
			}
		}
	}

	// 批量获取用户信息
	userInfoMap := make(map[string]types.Sender)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err == nil {
			for userId, userInfo := range userListResp.UserInfo {
				userInfoMap[userId] = types.Sender{
					UserID:   userId,
					NickName: userInfo.NickName,
					Avatar:   userInfo.Avatar,
					UserType: int8(userInfo.UserType),
				}
			}
		}
	}

	var list []types.Message
	for _, m := range detail.Content {
		var tMsg types.Message
		// 转换基本字段
		tMsg.Id = m.Id
		tMsg.MessageID = m.MessageID
		tMsg.ConversationID = m.ConversationID
		tMsg.ConversationType = m.ConversationType
		tMsg.CreatedAt = m.CreatedAt.String()
		tMsg.Seq = m.Seq

		// 转换消息内容
		if m.Msg != nil {
			msgJSON, _ := json.Marshal(m.Msg)
			json.Unmarshal(msgJSON, &tMsg.Msg)
		}

		// 填充发送者信息
		sendUserID := ""
		if m.SendUserID != nil {
			sendUserID = *m.SendUserID
		}

		if sendUserID != "" {
			if sender, exists := userInfoMap[sendUserID]; exists {
				tMsg.Sender = sender
			} else {
				tMsg.Sender = types.Sender{
					UserID:   sendUserID,
					NickName: "用户" + sendUserID[len(sendUserID)-4:], // 后四位辅助识别
					Avatar:   "",
				}
			}
		} else {
			tMsg.Sender = types.Sender{
				UserID:   "",
				NickName: "系统消息",
				Avatar:   "",
			}
		}

		list = append(list, tMsg)
	}

	return &types.GetForwardDetailsRes{
		List: list,
	}, nil
}
