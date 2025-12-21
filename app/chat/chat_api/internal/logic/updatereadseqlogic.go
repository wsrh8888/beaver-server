package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateReadSeqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReadSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReadSeqLogic {
	return &UpdateReadSeqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReadSeqLogic) UpdateReadSeq(req *types.UpdateReadSeqReq) (*types.UpdateReadSeqRes, error) {
	// 参数验证
	if req.ConversationID == "" {
		return nil, errors.New("ConversationID不能为空")
	}
	if req.ReadSeq < 0 {
		return nil, errors.New("ReadSeq值不对")
	}

	// 查询用户会话关系
	var userConvo chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, req.UserID).First(&userConvo).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果记录不存在，创建新记录
			version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)
			userConvo = chat_models.ChatUserConversation{
				UserID:         req.UserID,
				ConversationID: req.ConversationID,
				IsHidden:       false,
				IsPinned:       false,
				IsMuted:        false,
				UserReadSeq:    req.ReadSeq,
				Version:        version,
			}
			if err := l.svcCtx.DB.Create(&userConvo).Error; err != nil {
				l.Errorf("创建用户会话关系失败: userId=%s, conversationId=%s, error=%v", req.UserID, req.ConversationID, err)
				return nil, err
			}
			l.Infof("创建用户会话关系并设置已读序列号: userId=%s, conversationId=%s, readSeq=%d", req.UserID, req.ConversationID, req.ReadSeq)
		} else {
			l.Errorf("查询用户会话关系失败: userId=%s, conversationId=%s, error=%v", req.UserID, req.ConversationID, err)
			return nil, err
		}
	} else {
		// 如果记录存在，更新已读序列号（只有当新的readSeq大于当前值时才更新，避免回退）
		if req.ReadSeq > userConvo.UserReadSeq {
			version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)
			err = l.svcCtx.DB.Model(&userConvo).
				Updates(map[string]interface{}{
					"user_read_seq": req.ReadSeq,
					"updated_at":    time.Now(),
					"version":       version,
				}).Error
			if err != nil {
				l.Errorf("更新已读序列号失败: userId=%s, conversationId=%s, readSeq=%d, error=%v", req.UserID, req.ConversationID, req.ReadSeq, err)
				return nil, err
			}
			l.Infof("更新已读序列号成功: userId=%s, conversationId=%s, readSeq=%d (原值: %d)", req.UserID, req.ConversationID, req.ReadSeq, userConvo.UserReadSeq)

			// 异步推送WebSocket通知
			go l.notifyReadSeqUpdate(req.UserID, req.ConversationID, version)
		} else {
			l.Infof("已读序列号无需更新: userId=%s, conversationId=%s, 当前readSeq=%d, 请求readSeq=%d", req.UserID, req.ConversationID, userConvo.UserReadSeq, req.ReadSeq)
		}
	}

	return &types.UpdateReadSeqRes{
		Success: true,
	}, nil
}

// notifyReadSeqUpdate 推送已读序列号更新通知
func (l *UpdateReadSeqLogic) notifyReadSeqUpdate(userID, conversationID string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.Errorf("推送已读序列号更新通知时发生panic: %v", r)
		}
	}()

	// 构建 chat_user_conversations 表更新数据（WebSocket推送使用前端期望的表名）
	userConversationsUpdate := map[string]interface{}{
		"table":          "user_conversations", // 前端期望的表名
		"userId":         userID,
		"conversationId": conversationID,
		"data": []map[string]interface{}{
			{
				"version": version,
			},
		},
	}

	// 推送WebSocket通知（给自己，同步到其他设备）
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.CHAT_MESSAGE, wsTypeConst.ChatUserConversationReceive, userID, userID, map[string]interface{}{
		"tableUpdates": []map[string]interface{}{userConversationsUpdate},
	}, conversationID)

	l.Infof("推送已读序列号更新通知完成: userId=%s, conversationId=%s, version=%d", userID, conversationID, version)
}
