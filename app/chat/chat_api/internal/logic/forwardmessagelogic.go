package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/models/ctype"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ForwardMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForwardMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForwardMessageLogic {
	return &ForwardMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForwardMessageLogic) ForwardMessage(req *types.ForwardMessageReq) (resp *types.ForwardMessageRes, err error) {
	// 1. 获取原始消息对象列表
	var originMessages []chat_models.ChatMessage
	err = l.svcCtx.DB.Where("message_id IN ?", req.MessageIDs).Order("created_at asc").Find(&originMessages).Error
	if err != nil {
		l.Logger.Errorf("获取待转发消息失败: %v", err)
		return nil, err
	}

	if len(originMessages) == 0 {
		return nil, errors.New("未找到有效的转发消息")
	}

	if req.ForwardMode == 1 {
		// --- 逐条转发模式 ---
		for _, m := range originMessages {
			// 生成全新的客户端ID（简单起见使用UUID+原始ID，大厂通常由前端传入或后端补齐）
			newMsgID := uuid.New().String()

			// 调用 RPC 发送消息 (复用发送逻辑)
			_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
				UserId:         req.UserID,
				ConversationId: req.TargetID,
				MessageId:      newMsgID,
				Msg:            l.convertModelToProtoMsg(m.Msg),
			})
			if err != nil {
				l.Logger.Errorf("逐条转发失败: %v", err)
				// 商业化项目通常会继续处理下一条，或者返回部分成功的提示
			}
		}
	} else {
		// --- 合并转发模式 ---
		recordID := uuid.New().String()

		// 2. 将消息快照存入详情表 (冷数据)
		err = l.svcCtx.DB.Create(&chat_models.ChatForward{
			RecordID: recordID,
			Content:  originMessages, // 直接赋值，由 ForwardContent.Value 接口处理序列化
		}).Error
		if err != nil {
			l.Logger.Errorf("创建转发详情失败: %v", err)
			return nil, err
		}

		// 3. 发送合并转发卡片 (热数据)
		title := "聊天记录"
		if len(originMessages) > 0 {
			title = "群聊的聊天记录"
		}

		_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
			UserId:         req.UserID,
			ConversationId: req.TargetID,
			MessageId:      uuid.New().String(),
			Msg: &chat_rpc.Msg{
				Type: uint32(ctype.ForwardMsgType),
				ForwardMsg: &chat_rpc.ForwardMsg{
					Title:    title,
					RecordId: recordID,
					Count:    int32(len(originMessages)),
				},
			},
		})
		if err != nil {
			l.Logger.Errorf("发送合并转发卡片失败: %v", err)
			return nil, err
		}
	}

	return &types.ForwardMessageRes{
		ForwardTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// 辅助方法：将 DB Model 的 Msg 转换为 RPC 的 Msg (需要适配具体字段)
func (l *ForwardMessageLogic) convertModelToProtoMsg(m *ctype.Msg) *chat_rpc.Msg {
	if m == nil {
		return nil
	}
	// 这里最简单的方式是 JSON 中转，或者手动构建
	jsonData, _ := json.Marshal(m)
	var protoMsg chat_rpc.Msg
	json.Unmarshal(jsonData, &protoMsg)
	return &protoMsg
}
